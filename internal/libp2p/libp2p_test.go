package libp2p

import (
	"context"
	"testing"

	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/multiformats/go-multiaddr"
)

func Test_SetupLibp2p(t *testing.T) {
	pk, _, err := crypto.GenerateKeyPair(crypto.ECDSA, 2048)
	if err != nil {
		t.Fatal(err)
	}
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ps := pstoremem.NewPeerstore()
	ds := dssync.MutexWrap(datastore.NewMapDatastore())
	h, d, err := SetupLibp2p(ctx, pk, nil, []multiaddr.Multiaddr{addr}, ps, ds)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	if err := h.Close(); err != nil {
		t.Fatal(err)
	}
	if peers := DefaultBootstrapPeers(); len(peers) == 0 {
		t.Fatal("invalid peers received")
	}
}
