package rtns

import (
	"context"
	"testing"

	lp "github.com/RTradeLtd/rtns/internal/libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/multiformats/go-multiaddr"
)

var (
	ipfsPath = "/ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
)

func Test_New_Publisher(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pk := newPK(t)
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")
	if err != nil {
		t.Fatal(err)
	}
	publisher, err := NewPublisher(ctx, "test", pk, []multiaddr.Multiaddr{addr})
	if err != nil {
		t.Fatal(err)
	}
	publisher.Bootstrap(lp.DefaultBootstrapPeers())
	publisher.SetNameSys()
	if err := publisher.Publish(ctx, newPK(t), ipfsPath); err != nil {
		t.Fatal(err)
	}
	publisher.Close()
}

func newPK(t *testing.T) crypto.PrivKey {
	pk, _, err := crypto.GenerateKeyPair(crypto.ECDSA, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return pk
}
