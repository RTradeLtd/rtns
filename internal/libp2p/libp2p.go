package libp2p

import (
	"context"

	datastore "github.com/ipfs/go-datastore"
	config "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipns"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	libcore "github.com/libp2p/go-libp2p-core"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-core/peerstore"
	ipnet "github.com/libp2p/go-libp2p-core/pnet"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtOpts "github.com/libp2p/go-libp2p-kad-dht/opts"
	pnet "github.com/libp2p/go-libp2p-pnet"
	record "github.com/libp2p/go-libp2p-record"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/multiformats/go-multiaddr"
)

// temporal bootstrap nodes
var temporalBootstrap = []string{
	"/ip4/172.218.49.115/tcp/4002/ipfs/QmPvnFXWAz1eSghXD6JKpHxaGjbVo4VhBXY2wdBxKPbne5",
	"/ip4/172.218.49.115/tcp/4003/ipfs/QmXow5Vu8YXqvabkptQ7HddvNPpbLhXzmmU53yPCM54EQa",
	"/ip4/35.203.44.77/tcp/4001/ipfs/QmUMtzoRfQ6FttA7RygL8jJf7TZJBbdbZqKTmHfU6QC5Jm",
}

// SetupLibp2p returns a routed host and DHT instances that can be used to
// easily create a ipfslite Peer. The DHT is NOT bootstrapped. You may consider
// to use Peer.Bootstrap() after creating the IPFS-Lite Peer.
func SetupLibp2p(
	ctx context.Context,
	hostKey crypto.PrivKey,
	secret []byte,
	listenAddrs []multiaddr.Multiaddr,
	pstore peerstore.Peerstore,
	dstore datastore.Batching,
) (host.Host, *dht.IpfsDHT, error) {

	var (
		prot ipnet.Protector
		opts []libp2p.Option
		err  error
	)

	// Create protector if we have a secret.
	if secret != nil && len(secret) > 0 {
		var key [32]byte
		copy(key[:], secret)
		prot, err = pnet.NewV1ProtectorFromBytes(&key)
		if err != nil {
			return nil, nil, err
		}
		opts = append(opts, libp2p.PrivateNetwork(prot))
	}
	opts = append(opts,
		libp2p.Identity(hostKey),
		libp2p.ListenAddrs(listenAddrs...),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.DefaultMuxers,
		libp2p.DefaultTransports,
		libp2p.DefaultSecurity,
	)
	h, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, nil, err
	}

	idht, err := dht.New(ctx, h,
		dhtOpts.Validator(record.NamespacedValidator{
			"pk":   record.PublicKeyValidator{},
			"ipns": ipns.Validator{KeyBook: pstore},
		}),
		dhtOpts.Datastore(dstore),
	)
	if err != nil {
		h.Close()
		return nil, nil, err
	}
	rHost := routedhost.Wrap(h, idht)
	return rHost, idht, nil
}

// DefaultBootstrapPeers returns the default lsit
// of bootstrap peers, with added temporal production
// peers
func DefaultBootstrapPeers() []libcore.PeerAddrInfo {
	// conversion copied from go-ipfs
	defaults, _ := config.DefaultBootstrapPeers()
	tPeers, _ := config.ParseBootstrapPeers(temporalBootstrap)
	defaults = append(defaults, tPeers...)
	pinfos := make(map[peer.ID]*libcore.PeerAddrInfo)
	for _, bootstrap := range defaults {
		pinfo, ok := pinfos[bootstrap.ID]
		if !ok {
			pinfo = new(libcore.PeerAddrInfo)
			pinfos[bootstrap.ID] = pinfo
			pinfo.ID = bootstrap.ID
		}

		pinfo.Addrs = append(pinfo.Addrs, bootstrap.Addrs...)
	}

	var peers []libcore.PeerAddrInfo
	for _, pinfo := range pinfos {
		peers = append(peers, *pinfo)
	}
	return peers
}
