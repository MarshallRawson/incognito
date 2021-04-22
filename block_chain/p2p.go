package block_chain

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

type p2pStuff struct {
	peer_chan    chan peer.AddrInfo
	host         host.Host
	ctx          context.Context
	sub          *pubsub.Subscription
	pubsub       *pubsub.PubSub
	rendezvous   string
	topic        *pubsub.Topic
	genesis_hash [HashSize]byte
}

type Message string

const (
	publishBlock = "PublishBlock"
	requestChain = "RequestChain"
	publishChain = "PublishChain"
)

func (bc *BlockChain) Join(rendezvous string, genesis_hash [HashSize]byte) {
	bc.p2p = &p2pStuff{
		rendezvous:   rendezvous,
		genesis_hash: genesis_hash,
	}

	bc.p2p.ctx = context.Background()
	var err error
	bc.p2p.host, err = libp2p.New(bc.p2p.ctx,
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Identity(bc.self.priv))
	if err != nil {
		panic(err)
	}

	bc.p2p.pubsub, err = pubsub.NewGossipSub(bc.p2p.ctx, bc.p2p.host)
	if err != nil {
		panic(err)
	}
	ser, err := discovery.NewMdnsService(bc.p2p.ctx,
		bc.p2p.host,
		time.Hour,
		"incognito")
	if err != nil {
		panic(err)
	}
	ser.RegisterNotifee(bc)
	bc.p2p.topic, err = bc.p2p.pubsub.Join(rendezvous)
	if err != nil {
		panic(err)
	}

	bc.p2p.sub, err = bc.p2p.topic.Subscribe()
	if err != nil {
		panic(err)
	}

	go bc.readLoop()
}

func (p2p *p2pStuff) requestChain() {
	rc, err := json.Marshal([][]byte{[]byte(requestChain)})
	if err != nil {
		panic(err)
	}
	p2p.topic.Publish(p2p.ctx, rc)
}

func (bc *BlockChain) readLoop() {
	for {
		msg, err := bc.p2p.sub.Next(bc.p2p.ctx)
		if err != nil {
			fmt.Println("Error reading from libp2p peer: ", err)
		}
		if msg.ReceivedFrom == bc.self.iD {
			continue
		}
		var d [][]byte
		err = json.Unmarshal(msg.Data, &d)
		if err != nil {
			fmt.Println("Could not json unmarshall message: ", err)
		}
		switch Message(string(d[0])) {
		case publishBlock:
			b, err := UnmarshalBlock(d[1])
			if err != nil {
				fmt.Println("Could not Unmarshal block from message: ", err)
			}
			err = bc.AddBlock(b)
			if err != nil {
				bc.p2p.requestChain()
			}
		case requestChain:
			bc.p2p.publishChain(bc.chain)
		case publishChain:
			if len(d)-1 > bc.chain.Len() {
				// unmarshall the first (and check that it is a genesis)
				b, err := UnmarshalBlock(d[1])
				if err == nil && b.GetAction() == genesis && b.GetHash() == bc.p2p.genesis_hash {
					bc_test := New(MakeSelf(b.GetName()), false)
					chain_valid := true
					for i := 1; i < len(d); i += 1 {
						b, err = UnmarshalBlock(d[i])
						if err != nil {
							fmt.Println("Could not Unmarshal block from message: ", err)
							chain_valid = false
							break
						}
						err := bc_test.AddBlock(b)
						if err != nil {
							fmt.Println("invalid block: ", b.AsString())
							chain_valid = false
							break
						}
					}
					if chain_valid == true {
						bc.chain.Init()
						for i := 1; i < len(d); i += 1 {
							b, err = UnmarshalBlock(d[i])
							if err != nil {
								fmt.Println("Could not Unmarshal block from message: ", err)
								break
							}
							err := bc.AddBlock(b)
							if err != nil {
								fmt.Println("invalid block that is supposed to be valid: ", b.AsString())
							}
						}
					}
				}
			}
		}
	}
}

func (p2p *p2pStuff) publishBlock(b Block) {
	block_bytes, err := MarshalBlock(b)
	if err != nil {
		panic(err)
	}
	payload, err := json.Marshal([][]byte{[]byte(publishBlock), block_bytes})
	if err != nil {
		panic(err)
	}
	p2p.topic.Publish(p2p.ctx, payload)
}

func (p2p *p2pStuff) publishChain(chain list.List) {
	d := make([][]byte, chain.Len()+1)
	d[0] = []byte(publishChain)
	i := 1
	for b := chain.Front(); b != nil; b = b.Next() {
		m_b, err := MarshalBlock(b.Value.(Block))
		if err != nil {
			panic(err)
		}
		d[i] = m_b
		i += 1
	}
	payload, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	p2p.topic.Publish(p2p.ctx, payload)
}

func (bc *BlockChain) HandlePeerFound(pi peer.AddrInfo) {
	if bc.nodes.Has(pi.ID) || bc.nodes.Len() == 0 {
		err := bc.p2p.host.Connect(context.Background(), pi)
		if err != nil {
			fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
		}
	}
}
