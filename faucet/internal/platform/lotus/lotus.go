package lotus

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/v1api"
)

type API interface {
	NodeStatus(ctx context.Context, inclChainStatus bool) (api.NodeStatus, error)
	NetPeers(context.Context) ([]peer.AddrInfo, error)
	Version(context.Context) (api.APIVersion, error)
	ID(context.Context) (peer.ID, error)
}

func GetToken() (string, error) {
	lotusPath := os.Getenv("LOTUS_PATH")
	fmt.Println("LOTUS_PATH=", lotusPath)
	if lotusPath == "" {
		return "", fmt.Errorf("LOTUS_PATH not set in environment")
	}
	token, err := os.ReadFile(path.Join(lotusPath, "/token"))
	return string(token), err
}

func VerifyWallet(ctx context.Context, api v1api.FullNode, addr address.Address) error {
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
