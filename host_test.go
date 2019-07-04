package rtns

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	tutil "github.com/RTradeLtd/go-libp2p-testutils"
	pb "github.com/RTradeLtd/grpc/krab"
	kaas "github.com/RTradeLtd/kaas/v2"
	"github.com/RTradeLtd/rtns/mocks"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

var (
	ipfsPath1 = "/ipfs/QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
	ipfsPath2 = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
)

func Test_RTNS(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//////////////////
	// setup mocks //
	////////////////

	pk1 := tutil.NewPrivateKey(t)
	pk1Bytes, err := pk1.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	pk2 := tutil.NewPrivateKey(t)
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
	fns.PublishReturnsOnCall(2, errors.New("publish failed"))

	//////////////////////
	// setup publisher //
	////////////////////

	RTNS := newTestRTNS(ctx, t, fkb, fns)

	//////////////////
	// start tests //
	////////////////

	// ensure no previous records have been published
	if err := RTNS.republishEntries(); err != errNoRecordsPublisher {
		t.Fatal("wrong error received")
	}

	if err := RTNS.Publish(ctx, pk1, true, "pk1", ipfsPath1); err != nil {
		t.Fatal(err)
	}
	if len(RTNS.cache.list()) != 1 {
		fmt.Println("cache length:", len(RTNS.cache.list()))
		t.Fatal("invalid cache length")
	}
	pid, err := peer.IDFromPublicKey(pk1.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pk1", pid.String())

	if err := RTNS.Publish(ctx, pk2, true, "pk2", ipfsPath2); err != nil {
		t.Fatal(err)
	}
	if len(RTNS.cache.list()) != 2 {
		fmt.Println("cache length:", len(RTNS.cache.list()))
		t.Fatal("invalid cache length")
	}
	pid, err = peer.IDFromPublicKey(pk2.GetPublic())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pk2", pid.String())

	if err := RTNS.republishEntries(); err != nil {
		t.Fatal(err)
	}

	if err := RTNS.Publish(ctx, pk2, true, "pk2", ipfsPath2); err == nil {
		t.Fatal("error expected")
	}
}

func Test_Keystore(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fkb := &mocks.FakeServiceClient{}
	pk := tutil.NewPrivateKey(t)
	pkBytes, err := pk.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	fkb.HasPrivateKeyReturnsOnCall(0, &pb.Response{Status: "OK"}, nil)
	fkb.HasPrivateKeyReturnsOnCall(1, &pb.Response{Status: "BAD"}, errors.New("no key"))
	fkb.GetPrivateKeyReturnsOnCall(0, &pb.Response{Status: "OK", PrivateKey: pkBytes}, nil)
	fkb.GetPrivateKeyReturnsOnCall(1, &pb.Response{Status: "BAD"}, errors.New("no keys"))

	rk := newRKeystore(ctx, &kaas.Client{ServiceClient: fkb})

	// test has
	if exists, err := rk.Has("hello"); err != nil {
		t.Fatal(err)
	} else if !exists {
		t.Fatal("key should exist")
	}
	if exists, err := rk.Has("world"); err == nil {
		t.Fatal("error expected")
	} else if exists {
		t.Fatal("key should not exist")
	}

	// test put
	if err := rk.Put("abc", nil); err == nil {
		t.Fatal("error expected")
	}

	// test get
	if pkRet, err := rk.Get("abc"); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(pk, pkRet) {
		t.Fatal("keys should be equal")
	}
	if pkRet, err := rk.Get("abc"); err == nil {
		t.Fatal("error expected")
	} else if pkRet != nil {
		t.Fatal("pk should be nil")
	}

	// test delete
	if err := rk.Delete("abc"); err == nil {
		t.Fatal("error expected")
	}

	// test list
	if ids, err := rk.List(); err == nil {
		t.Fatal("error expected")
	} else if len(ids) != 0 {
		t.Fatal("bad key length returned")
	}
}

func newTestRTNS(ctx context.Context, t *testing.T, fkb *mocks.FakeServiceClient, fns *mocks.FakeNameSystem) *RTNS {
	ds := tutil.NewDatastore(t)
	ps := tutil.NewPeerstore(t)
	addrs := []multiaddr.Multiaddr{tutil.NewMultiaddr(t)}
	pk := tutil.NewPrivateKey(t)
	logger := tutil.NewLogger(t)
	keys := newRKeystore(ctx, &kaas.Client{ServiceClient: fkb})
	_, dht := tutil.NewLibp2pHostAndDHT(ctx, t, logger.Desugar(), ds, ps, pk, addrs, nil)
	rtns := NewRTNS(ctx, dht, ds, keys, 128)
	rtns.ns = fns
	return rtns
}
