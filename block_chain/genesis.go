package block_chain

import (
	"fmt"
	"github.com/golang-collections/collections/set"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Genesis struct {
	block
	AdminValid PrivValidation
	Title      string
	NodeID     peer.ID
}

func NewGenesis(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	admin_valid PrivValidation,
	title string,
	node_id peer.ID) *Genesis {

	gen := Genesis{
		AdminValid: admin_valid,
		Title:      title,
		NodeID:     node_id,
	}
	gen.PrevHash = prev_hash
	gen.Name = name
	gen.PublishValid = pub_valid
	gen.Action = genesis
	gen.Hash = [HashSize]byte{0}
	gen.Hash = Hash(&gen)
	return &gen
}

func (g *Genesis) GetAction() Action {
	return g.block.GetAction()
}

func (g *Genesis) AsString() string {
	return fmt.Sprintf("%s: Welcome to %s\n", g.Name, g.Title)
}

func (g *Genesis) GetHash() [HashSize]byte {
	return g.block.GetHash()
}

func (g *Genesis) GetPrevHash() [HashSize]byte {
	return g.block.GetPrevHash()
}

func (g *Genesis) GetName() string {
	return g.block.GetName()
}

func (g *Genesis) SetHash(new_hash [HashSize]byte) {
	g.block.SetHash(new_hash)
}

func (g *Genesis) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	return true
}

func (g *Genesis) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	publishers[g.PublishValid.NextPuzzle] = g.Name
	admins.Insert(g.AdminValid.NextPuzzle)
}
