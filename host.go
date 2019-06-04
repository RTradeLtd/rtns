package rtns

import (
	"context"
	"fmt"
	"time"

	kaas "github.com/RTradeLtd/kaas/v2"
	lp "github.com/RTradeLtd/rtns/internal/libp2p"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-path"
	pinfo "github.com/libp2p/go-libp2p-core"
	ci "github.com/libp2p/go-libp2p-core/crypto"
	host "github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peerstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/multiformats/go-multiaddr"
)

// Service implements the rtns logic
// as a consumable service or library.
type Service interface {
	// Close is used to trigger a shutdown of internal services
	Close()

	// DefaultBootstrapPeers returns the normal libp2p bootstrap peers, as well as the production nodes of Temporal.
	DefaultBootstrapPeers() []pinfo.PeerAddrInfo
	// Bootstrap is used to bootstrap the dht
	Bootstrap(peers []pinfo.PeerAddrInfo)

	// Publish enables publishing of an IPNS record with a default lifetime of 24 hours
	Publish(ctx context.Context, pk ci.PrivKey, cache bool, keyID, content string) error
	// PublishWithEOL allows specifying a lifetime for this record overriding the default lifetime of 24 hours
	PublishWithEOL(ctx context.Context, pk ci.PrivKey, eol time.Time, cache bool, keyID, content string) error

	// GetKey returns a key from the underlying krab keystore
	GetKey(name string) (ci.PrivKey, error)
	// HasKey returns whether or not the key is in our keystore
	HasKey(name string) (bool, error)
}

// Config is used to configure the RTNS service
type Config struct {
	Datastore   datastore.Batching
	PK          ci.PrivKey
	ListenAddrs []multiaddr.Multiaddr
	Secret      []byte
}

// rtns manages all the needed components
// to interact with public and private IPNS
// networks.
type rtns struct {
	h     host.Host
	pk    ci.PrivKey
	d     *dht.IpfsDHT
	ds    datastore.Datastore
	ns    namesys.NameSystem
	ps    peerstore.Peerstore
	ctx   context.Context
	keys  *rkeystore
	cache *cache
}

// NewService is used to instantiate an RTNS publisher service
func NewService(ctx context.Context, kbClient *kaas.Client, cfg Config) (Service, error) {
	return newRTNS(ctx, kbClient, cfg)
}

// NewRTNS is used to instantiate our RTNS service, and start the republisher
// intended to be used as a `Publisher` type by external libraries
func newRTNS(ctx context.Context, kbClient *kaas.Client, cfg Config) (*rtns, error) {
	ps := pstoremem.NewPeerstore()
	ht, dt, err := lp.SetupLibp2p(ctx, cfg.PK, cfg.Secret, cfg.ListenAddrs, ps, cfg.Datastore)
	if err != nil {
		return nil, err
	}
	r := &rtns{
		h:     ht,
		d:     dt,
		pk:    cfg.PK,
		ds:    cfg.Datastore,
		ps:    ps,
		ns:    namesys.NewNameSystem(dt, cfg.Datastore, 128),
		ctx:   ctx,
		keys:  newRKeystore(ctx, kbClient),
		cache: newCache(),
	}
	go r.startRepublisher()
	return r, nil
}

// Close is used to close all service needed by our publisher
func (r *rtns) Close() {
	if err := r.d.Close(); err != nil {
		fmt.Println("error shutting down dht:", err.Error())
	}
	if err := r.h.Close(); err != nil {
		fmt.Println("error shutting down host:", err.Error())
	}
	if err := r.ds.Close(); err != nil {
		fmt.Println("error shutting down datastore:", err.Error())
	}
}

// Publish enables publishing of an IPNS record with a default lifetime of 24 hours
func (r *rtns) Publish(ctx context.Context, pk ci.PrivKey, cache bool, keyID, content string) error {
	if cache {
		r.cache.Set(keyID)
	}
	return r.ns.Publish(ctx, pk, path.FromString(content))
}

// PublishWithEOL allows specifying a lifetime for this record overriding the default lifetime of 24 hours
func (r *rtns) PublishWithEOL(ctx context.Context, pk ci.PrivKey, eol time.Time, cache bool, keyID, content string) error {
	if cache {
		r.cache.Set(keyID)
	}
	return r.ns.PublishWithEOL(ctx, pk, path.FromString(content), eol)
}

// GetKey returns a key from the underlying krab keystore

func (r *rtns) GetKey(name string) (ci.PrivKey, error) {
	return r.keys.Get(name)
}

// HasKey returns whether or not the key is in our keystore
func (r *rtns) HasKey(name string) (bool, error) {
	return r.keys.Has(name)
}
