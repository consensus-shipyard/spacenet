package main

import (
	"context"
	"crypto/tls"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	datastore "github.com/ipfs/go-ds-leveldb"
	logging "github.com/ipfs/go-log/v2"
	"github.com/pkg/errors"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"
	"go.uber.org/zap"

	"github.com/filecoin-project/faucet/internal/faucet"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/api/v0api"

	app "github.com/filecoin-project/faucet/internal/http"
)

var build = "develop"

func main() {
	logger := logging.Logger("SPACENET-FAUCET")

	lvl, err := logging.LevelFromString("info")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)

	if err := run(logger); err != nil {
		logger.Fatalln("main: error:", err)
	}
}

func run(log *logging.ZapEventLogger) error {
	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:60s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			HTTPSHost       string        `conf:"default:0.0.0.0:443"`
			BackendHost     string        `conf:"required"`
			AllowedOrigins  []string      `conf:"required"`
		}
		TLS struct {
			CertFile string `conf:"required"`
			KeyFile  string `conf:"required"`
		}
		Filecoin struct {
			Address string `conf:"default:t1jlm55oqkdalh2l3akqfsaqmpjxgjd36pob34dqy"`
			// Amount of tokens that below is in FIL.
			TotalWithdrawalLimit   uint64 `conf:"default:10000"`
			AddressWithdrawalLimit uint64 `conf:"default:20"`
			WithdrawalAmount       uint64 `conf:"default:10"`
		}
		Lotus struct {
			APIHost   string `conf:"default:127.0.0.1:1230"`
			AuthToken string
		}
		DB struct {
			Path     string `conf:"default:./_db_data"`
			Readonly bool   `conf:"default:false"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "Spacenet Faucet Service",
		},
	}

	const prefix = "FAUCET"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return err
	}

	// =========================================================================
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	expvar.NewString("build").Set(build)

	// =========================================================================
	// Database Support

	log.Infow("startup", "status", "initializing database support", "path", cfg.DB.Path)

	db, err := datastore.NewDatastore(cfg.DB.Path, &datastore.Options{
		Compression: ldbopts.NoCompression,
		NoSync:      false,
		Strict:      ldbopts.StrictAll,
		ReadOnly:    cfg.DB.Readonly,
	})
	if err != nil {
		return fmt.Errorf("couldn't initialize leveldb database: %w", err)
	}

	defer func() {
		log.Infow("shutdown", "status", "stopping leveldb")
		db.Close()
	}()

	// =========================================================================
	// Initialize authentication support

	log.Infow("startup", "status", "initializing authentication support")

	var authToken string

	if cfg.Lotus.AuthToken == "" {
		authToken, err = getToken()
		if err != nil {
			return fmt.Errorf("error getting authentication token: %w", err)
		}
	} else {
		authToken = cfg.Lotus.AuthToken
	}
	header := http.Header{"Authorization": []string{"Bearer " + authToken}}

	// =========================================================================
	// Start Lotus client

	faucetAddr, err := address.NewFromString(cfg.Filecoin.Address)
	if err != nil {
		return fmt.Errorf("failed to parse Faucet address: %w", err)
	}

	log.Infow("startup", "status", "initializing Lotus support", "host", cfg.Lotus.APIHost)

	lotusNode, lotusCloser, err := client.NewFullNodeRPCV0(context.Background(), "ws://"+cfg.Lotus.APIHost+"/rpc/v0", header)
	if err != nil {
		return fmt.Errorf("connecting to Lotus failed: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping Lotus client support")
		lotusCloser()
	}()

	log.Infow("Successfully connected to Lotus node")

	// sanity-check to see if the node owns the key.
	if err := verifyWallet(context.Background(), lotusNode, faucetAddr); err != nil {
		return fmt.Errorf("faucet wallet sanity-check failed: %w", err)
	}

	// =========================================================================
	// Start API Service

	log.Infow("startup", "status", "initializing HTTP API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cert, err := tls.LoadX509KeyPair(cfg.TLS.CertFile, cfg.TLS.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS key pair: %w", err)
	}
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}

	api := http.Server{
		TLSConfig: tlsConfig,
		Addr:      cfg.Web.HTTPSHost,
		Handler: app.Handler(log, lotusNode, db, shutdown, &faucet.Config{
			FaucetAddress:          faucetAddr,
			AllowedOrigins:         cfg.Web.AllowedOrigins,
			BackendAddress:         cfg.Web.BackendHost,
			TotalWithdrawalLimit:   cfg.Filecoin.TotalWithdrawalLimit,
			AddressWithdrawalLimit: cfg.Filecoin.AddressWithdrawalLimit,
			WithdrawalAmount:       cfg.Filecoin.WithdrawalAmount,
		}),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
	}()

	// =========================================================================
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}

func getToken() (string, error) {
	lotusPath := os.Getenv("LOTUS_PATH")
	if lotusPath == "" {
		return "", fmt.Errorf("LOTUS_PATH not set in environment")
	}
	token, err := os.ReadFile(path.Join(lotusPath, "/token"))
	return string(token), err
}

func verifyWallet(ctx context.Context, api v0api.FullNode, addr address.Address) error {
	l, err := api.WalletList(ctx)
	if err != nil {
		return err
	}

	for _, w := range l {
		if w == addr {
			return nil
		}
	}
	return fmt.Errorf("faucet wallet not owned by peer targeted by faucet server")
}
