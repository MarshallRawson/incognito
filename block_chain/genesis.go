package block_chain

import (
	"fmt"
	"github.com/golang-collections/collections/set"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Genesis struct {
	block
	action     Action
	adminValid PrivValidation
	title      string
	nodeID     peer.ID
}

func NewGenesis(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	admin_valid PrivValidation,
	title string,
	node_id peer.ID) *Genesis {

	gen := Genesis{
		action:     genesis,
		adminValid: admin_valid,
		title:      title,
		nodeID:     node_id,
	}
	gen.prevHash = prev_hash
	gen.name = name
	gen.publishValid = pub_valid
	gen.hash = [HashSize]byte{0}
	gen.hash = Hash(&gen)
	return &gen
}

func (g *Genesis) AsString() string {
	return fmt.Sprintf("%s: Welcome to %s\n", g.name, g.title)
}

func (g *Genesis) GetHash() [HashSize]byte {
	return g.block.GetHash()
}

func (g *Genesis) GetPrevHash() [HashSize]byte {
	return g.block.GetPrevHash()
}

func (g *Genesis) SetHash(new_hash [HashSize]byte) {
	g.block.SetHash(new_hash)
}

func (g *Genesis) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	return true
}

func (g *Genesis) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	publishers[g.publishValid.nextPuzzle] = g.name
	admins.Insert(g.adminValid.nextPuzzle)
}
