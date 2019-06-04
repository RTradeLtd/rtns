package rtns

import (
	"context"
	"errors"

	pb "github.com/RTradeLtd/grpc/krab"
	kaas "github.com/RTradeLtd/kaas/v2"
	keystore "github.com/ipfs/go-ipfs-keystore"
	ci "github.com/libp2p/go-libp2p-core/crypto"
)

// ensure rkeystore satisfies
// the keystore.Keystore interface
var _ keystore.Keystore = (*rkeystore)(nil)

// rkeystore satisfies the keystore.Keystore
// interface, providing access to a kaas
// backend for secure key management
type rkeystore struct {
	kb  *kaas.Client
	ctx context.Context
}

// newRKeystore implements a keystore.Keystore
// compatible version of the kaas client
func newRKeystore(ctx context.Context, kb *kaas.Client) *rkeystore {
	return &rkeystore{kb, ctx}
}

// Has returns whether or not a key exist in the Keystore
func (rk *rkeystore) Has(name string) (bool, error) {
	_, err := rk.kb.HasPrivateKey(rk.ctx, &pb.KeyGet{Name: name})
	if err != nil {
		return false, err
	}
	return true, nil
}

// Put stores a key in the Keystore, if a key with the same name already exists, returns ErrKeyExists
func (rk *rkeystore) Put(name string, pk ci.PrivKey) error {
	return errors.New("key puts not permitted")
}

// Get retrieves a key from the Keystore if it exists, and returns ErrNoSuchKey
// otherwise.
func (rk *rkeystore) Get(name string) (ci.PrivKey, error) {
	resp, err := rk.kb.GetPrivateKey(rk.ctx, &pb.KeyGet{Name: name})
	if err != nil {
		return nil, err
	}
	return ci.UnmarshalPrivateKey(resp.GetPrivateKey())
}

// Delete removes a key from the Keystore
func (rk *rkeystore) Delete(string) error {
	return errors.New("key deletes not permitted")
}

// List returns a list of key identifier
func (rk *rkeystore) List() ([]string, error) {
	return nil, errors.New("list not implemented")
}
