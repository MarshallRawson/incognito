package block_chain

import (
	"fmt"
	"github.com/golang-collections/collections/set"
	"github.com/libp2p/go-libp2p-core/peer"
)

type AddNode struct {
	block
	AdminValid PrivValidation
	NewNodeID  peer.ID
}

func NewAddNode(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	admin_valid PrivValidation,
	new_node_id peer.ID) *AddNode {

	an := AddNode{
		AdminValid: admin_valid,
		NewNodeID:  new_node_id,
	}

	an.PrevHash = prev_hash
	an.Name = name
	an.PublishValid = pub_valid
	an.Action = addNode
	an.Hash = [HashSize]byte{0}
	an.Hash = Hash(&an)
	return &an
}

func (an *AddNode) GetAction() Action {
	return an.block.GetAction()
}

func (an *AddNode) AsString() string {
	return fmt.Sprintf("%s: Added Node: %s\n", an.Name, an.NewNodeID.Pretty())
}

func (an *AddNode) GetHash() [HashSize]byte {
	return an.block.GetHash()
}

func (an *AddNode) GetPrevHash() [HashSize]byte {
	return an.block.GetPrevHash()
}

func (an *AddNode) GetName() string {
	return an.block.GetName()
}

func (an *AddNode) SetHash(new_hash [HashSize]byte) {
	an.block.SetHash(new_hash)
}

func (an *AddNode) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	if an.block.CheckValidations(publishers) == false {
		return false
	}
	if publishers[Hash(an.PublishValid.Solution)] != an.Name {
		return false
	}
	if admins.Has(Hash(an.AdminValid.Solution)) == false {
		return false
	}
	return true
}

func (an *AddNode) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	name := publishers[Hash(an.PublishValid.Solution)]
	delete(publishers, Hash(an.PublishValid.Solution))
	publishers[an.PublishValid.NextPuzzle] = name

	admins.Remove(Hash(an.AdminValid.Solution))
	admins.Insert(an.AdminValid.NextPuzzle)
}
