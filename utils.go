package rtns

import (
	"context"
	"fmt"
	"sync"
	"time"

	lp "github.com/RTradeLtd/rtns/internal/libp2p"
	libp2p "github.com/libp2p/go-libp2p-core"
)

// DefaultBootstrapPeers returns the normal libp2p bootstrap peers,
// as well as the production nodes of Temporal.
func (r *rtns) DefaultBootstrapPeers() []libp2p.PeerAddrInfo {
	return lp.DefaultBootstrapPeers()
}

// Bootstrap is an optional helper to connect to the given peers and bootstrap
// the Peer DHT (and Bitswap). This is a best-effort function. Errors are only
// logged and a warning is printed when less than half of the given peers
// could be contacted. It is fine to pass a list where some peers will not be
// reachable.
func (r *rtns) Bootstrap(peers []libp2p.PeerAddrInfo) {
	// prevent hung processes, timeout this helper call after 2 minutes
	cctx, cancel := context.WithTimeout(r.ctx, time.Minute*2)
	defer cancel()
	connected := make(chan struct{})

	var wg sync.WaitGroup
	for _, pinfo := range peers {
		//h.Peerstore().AddAddrs(pinfo.ID, pinfo.Addrs, peerstore.PermanentAddrTTL)
		wg.Add(1)
		go func(pinfo libp2p.PeerAddrInfo) {
			defer wg.Done()
			err := r.h.Connect(cctx, pinfo)
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

	err := r.d.Bootstrap(cctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}
