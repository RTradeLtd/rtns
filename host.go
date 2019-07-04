package rtns

import (
	"context"
	"time"

	"github.com/ipfs/go-datastore"
	keystore "github.com/ipfs/go-ipfs-keystore"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-path"
	ci "github.com/libp2p/go-libp2p-core/crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

// RTNS manages all the needed components
// to interact with public and private IPNS
// networks.
type RTNS struct {
	ns    namesys.NameSystem
	ds    datastore.Datastore
	ctx   context.Context
	keys  keystore.Keystore
	cache *cache
}

// NewRTNS instantiates our RTNS service, and starts the republisher
func NewRTNS(ctx context.Context, dt *dht.IpfsDHT, ds datastore.Datastore, keys keystore.Keystore, size int) *RTNS {
	r := &RTNS{
		ns:    namesys.NewNameSystem(dt, ds, size),
		ds:    ds,
		keys:  keys,
		cache: newCache(),
	}
	go r.startRepublisher()
	return r
}

// Publish enables publishing of an IPNS record with a default lifetime of 24 hours
func (r *RTNS) Publish(ctx context.Context, pk ci.PrivKey, cache bool, keyID, content string) error {
	if cache {
		r.cache.Set(keyID)
	}
	return r.ns.Publish(ctx, pk, path.FromString(content))
}

// PublishWithEOL allows specifying a lifetime for this record overriding the default lifetime of 24 hours
func (r *RTNS) PublishWithEOL(ctx context.Context, pk ci.PrivKey, eol time.Time, cache bool, keyID, content string) error {
	if cache {
		r.cache.Set(keyID)
	}
	return r.ns.PublishWithEOL(ctx, pk, path.FromString(content), eol)
}

// GetKey returns a key from the underlying krab keystore
func (r *RTNS) GetKey(name string) (ci.PrivKey, error) {
	return r.keys.Get(name)
}

// HasKey returns whether or not the key is in our keystore
func (r *RTNS) HasKey(name string) (bool, error) {
	return r.keys.Has(name)
}
