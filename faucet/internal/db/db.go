package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-datastore"
	"github.com/pkg/errors"

	"github.com/filecoin-project/faucet/internal/data"
	"github.com/filecoin-project/go-address"
)

var (
	totalInfoKey = datastore.NewKey("total_info_key")
)

type Database struct {
	store datastore.Datastore
}

func NewDatabase(store datastore.Datastore) *Database {
	return &Database{
		store: store,
	}
}

func (db *Database) GetTotalInfo(ctx context.Context) (data.TotalInfo, error) {
	var info data.TotalInfo

	b, err := db.store.Get(ctx, totalInfoKey)
	if err != nil && !errors.Is(err, datastore.ErrNotFound) {
		return data.TotalInfo{}, fmt.Errorf("failed to get total info: %w", err)
	}
	if errors.Is(err, datastore.ErrNotFound) {
		return info, nil
	}
	if err := json.Unmarshal(b, &info); err != nil {
		return data.TotalInfo{}, fmt.Errorf("failed to decode total info: %w", err)
	}
	return info, nil
}

func (db *Database) GetAddrInfo(ctx context.Context, addr address.Address) (data.AddrInfo, error) {
	var info data.AddrInfo

	b, err := db.store.Get(ctx, addrKey(addr))
	if err != nil && !errors.Is(err, datastore.ErrNotFound) {
		return data.AddrInfo{}, fmt.Errorf("failed to get addr info: %w", err)
	}
	if errors.Is(err, datastore.ErrNotFound) {
		return info, nil
	}
	if err := json.Unmarshal(b, &info); err != nil {
		return data.AddrInfo{}, fmt.Errorf("failed to decode addr info: %w", err)
	}
	return info, nil
}

func (db *Database) UpdateAddrInfo(ctx context.Context, targetAddr address.Address, info data.AddrInfo) error {
	bytes, err := json.Marshal(info)
	if err != nil {
		return err
	}

	err = db.store.Put(ctx, addrKey(targetAddr), bytes)
	if err != nil {
		return fmt.Errorf("failed to put addr info into db: %w", err)
	}

	return nil
}

func (db *Database) UpdateTotalInfo(ctx context.Context, info data.TotalInfo) error {
	bytes, err := json.Marshal(info)
	if err != nil {
		return err
	}

	err = db.store.Put(ctx, totalInfoKey, bytes)
	if err != nil {
		return fmt.Errorf("failed to put total info into db: %w", err)
	}

	return nil
}

func addrKey(addr address.Address) datastore.Key {
	return datastore.NewKey(addr.String() + ":value")
}
