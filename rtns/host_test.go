package rtns

import (
	"context"
	"fmt"
	"testing"

	lp "github.com/RTradeLtd/rtns/internal/libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/multiformats/go-multiaddr"
)

var (
	ipfsPath1 = "/ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
	ipfsPath2 = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
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
	pk1 := newPK(t)
	pk2 := newPK(t)
	publisher.Bootstrap(lp.DefaultBootstrapPeers())
	if err := publisher.Publish(ctx, pk1, ipfsPath1); err != nil {
		t.Fatal(err)
	}
	pid, err := peer.IDFromPublicKey(pk1.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pk1", pid.String())
	if err := publisher.Publish(ctx, pk2, ipfsPath2); err != nil {
		t.Fatal(err)
	}
	pid, err = peer.IDFromPublicKey(pk2.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pk2", pid.String())
	publisher.Close()
}

func newPK(t *testing.T) crypto.PrivKey {
	pk, _, err := crypto.GenerateKeyPair(crypto.ECDSA, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return pk
}
