package rtns

import (
	"context"
	"fmt"
	"testing"

	cfg "github.com/RTradeLtd/config/v2"
	pb "github.com/RTradeLtd/grpc/krab"
	kaas "github.com/RTradeLtd/kaas/v2"
	lp "github.com/RTradeLtd/rtns/internal/libp2p"
	"github.com/RTradeLtd/rtns/mocks"
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

	//////////////////
	// setup mocks //
	////////////////

	pk1 := newPK(t)
	pk1Bytes, err := pk1.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	pk2 := newPK(t)
	pk2Bytes, err := pk2.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	fkb := &mocks.FakeServiceClient{}
	fkb.GetPrivateKeyReturnsOnCall(0, &pb.Response{Status: "Ok", PrivateKey: pk1Bytes}, nil)
	fkb.GetPrivateKeyReturnsOnCall(1, &pb.Response{Status: "Ok", PrivateKey: pk2Bytes}, nil)

	fns := &mocks.FakeNameSystem{}
	fns.PublishReturnsOnCall(0, nil)
	fns.PublishReturnsOnCall(1, nil)
	fns.PublishReturnsOnCall(2, nil)
	fns.PublishReturnsOnCall(3, nil)

	//////////////////////
	// setup publisher //
	////////////////////

	publisher := newPublisher(ctx, t, fkb, fns)
	defer publisher.Close()
	publisher.Bootstrap(lp.DefaultBootstrapPeers())

	//////////////////
	// start tests //
	////////////////

	// ensure no previous records have been published
	if err := publisher.republishEntries(); err != errNoRecordsPublisher {
		t.Fatal("wrong error received")
	}

	if err := publisher.Publish(ctx, pk1, "pk1", ipfsPath1); err != nil {
		t.Fatal(err)
	}
	if len(publisher.cache.List()) != 1 {
		fmt.Println("cache length:", len(publisher.cache.List()))
		t.Fatal("invalid cache length")
	}
	pid, err := peer.IDFromPublicKey(pk1.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pk1", pid.String())

	if err := publisher.Publish(ctx, pk2, "pk2", ipfsPath2); err != nil {
		t.Fatal(err)
	}
	if len(publisher.cache.List()) != 2 {
		fmt.Println("cache length:", len(publisher.cache.List()))
		t.Fatal("invalid cache length")
	}
	pid, err = peer.IDFromPublicKey(pk2.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pk2", pid.String())

	if err := publisher.republishEntries(); err != nil {
		t.Fatal(err)
	}
}

func newPublisher(ctx context.Context, t *testing.T, fkb *mocks.FakeServiceClient, fns *mocks.FakeNameSystem) *Publisher {
	pk := newPK(t)
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")
	if err != nil {
		t.Fatal(err)
	}
	publisher, err := NewPublisher(ctx, newServicesConfig(), "test", pk, []multiaddr.Multiaddr{addr})
	if err != nil {
		t.Fatal(err)
	}
	publisher.keys.kb = &kaas.Client{ServiceClient: fkb}
	publisher.ns = fns
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
