package rtns

import (
	kaas "github.com/RTradeLtd/kaas/v2"
	keystore "github.com/ipfs/go-ipfs-keystore"
	ci "github.com/libp2p/go-libp2p-crypto"
)

// this is a hacky work-around in satisfying
// the keystore.Keystore interface with a gRPC backend

var _ keystore.Keystore = (*RKeystore)(nil)

// RKeystore satisfies the keystore.Keystore
// interface, providing access to a kaas
// backend for secure key management
type RKeystore struct {
	kb *kaas.Client
}

// NewRKeystore implements a keystore.Keystore
// compatible version of the kaas client
func NewRKeystore(kb *kaas.Client) *RKeystore {
	return &RKeystore{kb}
}

// Has returns whether or not a key exist in the Keystore
func (rk *RKeystore) Has(string) (bool, error) {
	return false, nil
}

// Put stores a key in the Keystore, if a key with the same name already exists, returns ErrKeyExists
func (rk *RKeystore) Put(string, ci.PrivKey) error {
	return nil
}

// Get retrieves a key from the Keystore if it exists, and returns ErrNoSuchKey
// otherwise.
func (rk *RKeystore) Get(string) (ci.PrivKey, error) {
	return nil, nil
}

// Delete removes a key from the Keystore
func (rk *RKeystore) Delete(string) error {
	return nil
}

// List returns a list of key identifier
func (rk *RKeystore) List() ([]string, error) {
	return nil, nil
}
