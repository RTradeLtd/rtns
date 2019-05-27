package rtns

import (
	"context"
	"errors"
	"fmt"
	"time"

	namesys "github.com/ipfs/go-ipfs/namesys"
	path "github.com/ipfs/go-path"

	proto "github.com/gogo/protobuf/proto"
	ds "github.com/ipfs/go-datastore"
	pb "github.com/ipfs/go-ipns/pb"
	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

///////////////////////////////////////////////////////////////////
// modified version of the one contained in ipfs/go-ipfs/namesys //
///////////////////////////////////////////////////////////////////

var errNoEntry = errors.New("no previous entry")

// DefaultRebroadcastInterval is the default interval at which we rebroadcast IPNS records
var DefaultRebroadcastInterval = time.Hour * 4

// FailureRetryInterval is the interval at which we retry IPNS records broadcasts (when they fail)
var FailureRetryInterval = time.Minute * 5

// DefaultRecordLifetime is the default lifetime for IPNS records
const DefaultRecordLifetime = time.Hour * 24

// StartRepublisher is used to start our republisher service
func (p *Publisher) StartRepublisher() {
	timer := time.NewTimer(DefaultRebroadcastInterval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if err := p.republishEntries(); err != nil {
				fmt.Println("failed to republish entries", err.Error())
				timer.Reset(FailureRetryInterval)
			}
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *Publisher) republishEntries() error {
	keys := p.cache.List()
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		priv, err := p.keys.Get(key)
		if err != nil {
			return errNoEntry
		}
		if err := p.republishEntry(p.ctx, priv); err != nil {
			return err
		}
	}
	return nil
}

func (p *Publisher) republishEntry(ctx context.Context, priv ci.PrivKey) error {
	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return err
	}

	// Look for it locally only
	lv, err := p.getLastVal(id)
	if err != nil {
		if err == errNoEntry {
			return nil
		}
		return err
	}

	// update record with same sequence number
	eol := time.Now().Add(DefaultRecordLifetime)
	return p.ns.PublishWithEOL(ctx, priv, lv, eol)
}

func (p *Publisher) getLastVal(id peer.ID) (path.Path, error) {
	// Look for it locally only
	val, err := p.ds.Get(namesys.IpnsDsKey(id))
	switch err {
	case nil:
	case ds.ErrNotFound:
		return "", errNoEntry
	default:
		return "", err
	}

	e := new(pb.IpnsEntry)
	if err := proto.Unmarshal(val, e); err != nil {
		return "", err
	}
	return path.Path(e.Value), nil
}
