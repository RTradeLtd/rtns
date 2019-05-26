package rtns

import (
	"context"
	"fmt"

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

// Publisher provides a helper to publish IPNS records
type Publisher struct {
	h   host.Host
	pk  ci.PrivKey
	d   *dht.IpfsDHT
	ds  datastore.Datastore
	ns  namesys.NameSystem
	ps  peerstore.Peerstore
	ctx context.Context
}

// NewPublisher is used to instantiate our IPNS publisher service
// NOTE: this DHT isn't bootstrapped
func NewPublisher(ctx context.Context, dsPath string, pk ci.PrivKey, listenAddrs []multiaddr.Multiaddr) (*Publisher, error) {
	ps := pstoremem.NewPeerstore()
	ds := dssync.MutexWrap(datastore.NewMapDatastore())
	ht, dt, err := lp.SetupLibp2p(ctx, pk, nil, listenAddrs, ps, ds)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		h:   ht,
		d:   dt,
		pk:  pk,
		ds:  ds,
		ps:  ps,
		ctx: ctx,
	}, nil
}

func (p *Publisher) SetNameSys() {
	p.ns = namesys.NewNameSystem(p.d, p.ds, 128)
}

// Close is used to close all service needed by our publisher
func (p *Publisher) Close() {
	if err := p.d.Close(); err != nil {
		fmt.Println("error shutting down dht:", err.Error())
	}
	if err := p.h.Close(); err != nil {
		fmt.Println("error shutting down host:", err.Error())
	}
	if err := p.ds.Close(); err != nil {
		fmt.Println("error shutting down datastore:", err.Error())
	}
}

// Publish is used to publish content with a fixed lifetime and ttl
func (p *Publisher) Publish(ctx context.Context, pk ci.PrivKey, content string) error {
	return p.ns.Publish(ctx, pk, path.FromString(content))
}
