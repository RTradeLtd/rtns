package rtns

import (
	"context"
	"fmt"
	"testing"

	cfg "github.com/RTradeLtd/config/v2"
	lp "github.com/RTradeLtd/rtns/internal/libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/multiformats/go-multiaddr"
)

var (
	ipfsPath1 = "/ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
	ipfsPath2 = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
)

// TODO:
// we need to configure fakes for the krab client
// so that we may spoof a valid krab backend

func Test_New_Publisher(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	publisher := newPublisher(ctx, t)
	defer publisher.Close()
	publisher.Bootstrap(lp.DefaultBootstrapPeers())

	// ensure no previous records have been published
	if err := publisher.republishEntries(); err == nil {
		t.Fatal("error expected")
	}

	pk1 := newPK(t)
	if err := publisher.Publish(ctx, pk1, "pk1", ipfsPath1); err != nil {
		t.Fatal(err)
	}
	if len(publisher.cache.List()) != 1 {
		t.Fatal("invalid cache length")
	}
	pid, err := peer.IDFromPublicKey(pk1.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pk1", pid.String())

	pk2 := newPK(t)
	if err := publisher.Publish(ctx, pk2, "pk1", ipfsPath2); err != nil {
		t.Fatal(err)
	}
	if len(publisher.cache.List()) != 2 {
		t.Fatal("invalid cache length")
	}
	pid, err = peer.IDFromPublicKey(pk2.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pk2", pid.String())

}

func newPublisher(ctx context.Context, t *testing.T) *Publisher {
	pk := newPK(t)
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")
	if err != nil {
		t.Fatal(err)
	}
	publisher, err := NewPublisher(ctx, newServicesConfig(), "test", pk, []multiaddr.Multiaddr{addr})
	if err != nil {
		t.Fatal(err)
	}
	return publisher
}
func newPK(t *testing.T) crypto.PrivKey {
	pk, _, err := crypto.GenerateKeyPair(crypto.ECDSA, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return pk
}

func newServicesConfig() cfg.Services {
	return cfg.Services{}
}
