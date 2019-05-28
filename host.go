package rtns

import (
	"context"
	"fmt"
	"time"

	cfg "github.com/RTradeLtd/config/v2"
	kaas "github.com/RTradeLtd/kaas/v2"
	lp "github.com/RTradeLtd/rtns/internal/libp2p"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/ipfs/go-ipfs/namesys"
	"github.com/ipfs/go-path"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/multiformats/go-multiaddr"
)

// RTNS is a standalone IPNS publishing service
// for use with the kaas keystore enabling secure
// management of IPNS records
type RTNS struct {
	h     host.Host
	pk    ci.PrivKey
	d     *dht.IpfsDHT
	ds    datastore.Datastore
	ns    namesys.NameSystem
	ps    peerstore.Peerstore
	ctx   context.Context
	keys  *RKeystore
	cache *Cache
}

// NewRTNS is used to instantiate our RTNS service
// NOTE: this DHT isn't bootstrapped
func NewRTNS(ctx context.Context, krabConfig cfg.Services, dsPath string, pk ci.PrivKey, listenAddrs []multiaddr.Multiaddr) (*RTNS, error) {
	ps := pstoremem.NewPeerstore()
	ds := dssync.MutexWrap(datastore.NewMapDatastore())
	ht, dt, err := lp.SetupLibp2p(ctx, pk, nil, listenAddrs, ps, ds)
	if err != nil {
		return nil, err
	}
	kb1, err := kaas.NewClient(krabConfig, false)
	if err != nil {
		return nil, err
	}
	r := &RTNS{
		h:     ht,
		d:     dt,
		pk:    pk,
		ds:    ds,
		ps:    ps,
		ns:    namesys.NewNameSystem(dt, ds, 128),
		ctx:   ctx,
		keys:  NewRKeystore(ctx, kb1),
		cache: NewCache(),
	}
	go r.startRepublisher()
	return r, nil
}

// Close is used to close all service needed by our publisher
func (r *RTNS) Close() {
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

// Publish is used to publish content with a fixed lifetime and ttl
func (r *RTNS) Publish(ctx context.Context, pk ci.PrivKey, cache bool, keyID, content string) error {
	if cache {
		r.cache.Set(keyID)
	}
	return r.ns.Publish(ctx, pk, path.FromString(content))
}

// PublishWithEOL is used to publish an IPNS record with non default lifetime values
func (r *RTNS) PublishWithEOL(ctx context.Context, pk ci.PrivKey, content string, eol time.Time) error {
	return r.ns.PublishWithEOL(ctx, pk, path.FromString(content), eol)
}
