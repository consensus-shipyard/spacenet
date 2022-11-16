package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/ipfs/go-datastore"
	levelds "github.com/ipfs/go-ds-leveldb"
	logging "github.com/ipfs/go-log/v2"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
)

var log = logging.Logger("faucet")

type server struct {
	lotus  v0api.FullNode
	db     datastore.Datastore
	faucet address.Address
}

type fundRequest struct {
	// Value abi.TokenAmount
	Address string
}

type fundResponse struct {
	Error string
}

func main() {
	// setting up logger
	lvl, err := logging.LevelFromString("info")
	if err != nil {
		panic(err)
	}
	logging.SetAllLoggers(lvl)

	// Starting lotus node.
	// run lotus auth create-token --perm <read,write,sign,admin>
	// to get an auth-token
	authToken, err := getToken()
	if err != nil {
		log.Errorf("error getting authentication token: %w", err)
		panic("couldn't get API token for lotus node")
	}
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	// FIXME: Pass this value in a command line flag.
	addr := "127.0.0.1:1230"

	// FIXME: Make this configurable
	db, err := NewLevelDB("./db", false)
	if err != nil {
		log.Errorf("couldn´t initialize leveldb database: %w", err)
	}

	res, closer, err := client.NewFullNodeRPCV0(context.Background(), "ws://"+addr+"/rpc/v0", headers)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()
	log.Infof("Successfully connected to Lotus node")
	s := server{lotus: res, db: db}

	// Starting http server.
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/fund", s.fundRequest).Methods("POST")
	r.HandleFunc("/", home)

	// Serve static files
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	// FIXME: Make this configurable
	port := ":8000"
	log.Infof("HTTP server listening at port: %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Errorf("error starting http server: %w", err)
		panic("error starting server")
	}

}

func getToken() (string, error) {
	lotusPath := os.Getenv("LOTUS_PATH")
	if lotusPath == "" {
		return "", fmt.Errorf("LOTUS_PATH not set in environment")
	}
	token, err := os.ReadFile(path.Join(lotusPath, "/token"))
	return string(token), err
}

// serves index file
func home(w http.ResponseWriter, r *http.Request) {
	p := path.Dir("./static/index.html")
	// set header
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}

// TODO: Finalize this method.
func (s *server) fundRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("error decoding request: %w", err)
		return
	}
	fmt.Println(">>>>>>>> Body", string(body))
	var req fundRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Errorf("error unmarshaling json request: %w", err)
		errResponse(w, "error unmarshaling request")
		return
	}

	addr, err := address.NewFromString(req.Address)
	if err != nil {
		errResponse(w, err.Error())
		return
	}
	err = s.fundAddr(addr, abi.NewTokenAmount(1000))
	if err != nil {
		errResponse(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	errResponse(w, "")
	return

}

func errResponse(w http.ResponseWriter, errStr string) {
	resp := fundResponse{Error: errStr}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func (s *server) fundAddr(addr address.Address, value abi.TokenAmount) error {
	// TODO: verify if the address is allowed to receive funds this soon.
	// (prevent DDoSing)
	ctx := context.TODO()

	FaucetAddr, err := address.NewFromString("f1cp4q4lqsdhob23ysywffg2tvbmar5cshia4rweq")
	if err != nil {
		panic(err)
	}
	msg, err := s.lotus.MpoolPushMessage(ctx, &types.Message{
		To:     addr,
		From:   FaucetAddr,
		Value:  value,
		Method: 0, // methodSend
		Params: nil,
	}, nil)
	if err != nil {
		log.Errorw("Error pushing join subnet message to parent api", "err", err)
		return err
	}

	// wait state message.
	_, err = s.lotus.StateWaitMsg(ctx, msg.Cid(), build.MessageConfidence)
	if err != nil {
		log.Errorw("Error waiting for message to be committed", "err", err)
		return err
	}

	log.Infow("Address funded successfully")
	return nil
}

func NewLevelDB(path string, readonly bool) (datastore.Batching, error) {
	return levelds.NewDatastore(path, &levelds.Options{
		Compression: ldbopts.NoCompression,
		NoSync:      false,
		Strict:      ldbopts.StrictAll,
		ReadOnly:    readonly,
	})
}
