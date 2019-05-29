package rtns

import (
	"fmt"
	"sync"

	lp "github.com/RTradeLtd/rtns/internal/libp2p"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

// DefaultBootstrap is a wrapper around Bootstrap
// that handles bootstrapping to the default libp2p bootstrap
// nodes, as well as the Temporal production nodes
func (r *rtns) DefaultBootstrap() {
	r.Bootstrap(lp.DefaultBootstrapPeers())
}

// Bootstrap is an optional helper to connect to the given peers and bootstrap
// the Peer DHT (and Bitswap). This is a best-effort function. Errors are only
// logged and a warning is printed when less than half of the given peers
// could be contacted. It is fine to pass a list where some peers will not be
// reachable.
func (r *rtns) Bootstrap(peers []peerstore.PeerInfo) {
	connected := make(chan struct{})

	var wg sync.WaitGroup
	for _, pinfo := range peers {
		//h.Peerstore().AddAddrs(pinfo.ID, pinfo.Addrs, peerstore.PermanentAddrTTL)
		wg.Add(1)
		go func(pinfo peerstore.PeerInfo) {
			defer wg.Done()
			err := r.h.Connect(r.ctx, pinfo)
			if err != nil {
				fmt.Println("error", err.Error())
				return
			}
			fmt.Println("Connected to", pinfo.ID)
			connected <- struct{}{}
		}(pinfo)
	}

	go func() {
		wg.Wait()
		close(connected)
	}()

	i := 0
	for range connected {
		i++
	}
	if nPeers := len(peers); i < nPeers/2 {
		fmt.Printf("only connected to %d bootstrap peers out of %d\n", i, nPeers)
	}

	err := r.d.Bootstrap(r.ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}
