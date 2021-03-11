package block_chain

import (
	"fmt"
	"github.com/golang-collections/collections/set"
	"github.com/libp2p/go-libp2p-core/peer"
)

type AddNode struct {
	block
	action     Action
	adminValid PrivValidation
	newNodeID  peer.ID
}

func NewAddNode(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	admin_valid PrivValidation,
	new_node_id peer.ID) *AddNode {

	an := AddNode{
		action:     addNode,
		adminValid: admin_valid,
		newNodeID:  new_node_id,
	}

	an.prevHash = prev_hash
	an.name = name
	an.publishValid = pub_valid
	an.hash = [HashSize]byte{0}
	an.hash = Hash(&an)
	return &an
}

func (an *AddNode) AsString() string {
	return fmt.Sprintf("%s: Added Node: %s\n", an.name, an.newNodeID.Pretty())
}

func (an *AddNode) GetHash() [HashSize]byte {
	return an.block.GetHash()
}

func (an *AddNode) GetPrevHash() [HashSize]byte {
	return an.block.GetPrevHash()
}

func (an *AddNode) SetHash(new_hash [HashSize]byte) {
	an.block.SetHash(new_hash)
}

func (an *AddNode) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	if an.block.CheckValidations(publishers) == false {
		return false
	}
	if publishers[Hash(an.publishValid.solution)] != an.name {
		return false
	}
	if admins.Has(Hash(an.adminValid.solution)) == false {
		return false
	}
	return true
}

func (an *AddNode) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	name := publishers[Hash(an.publishValid.solution)]
	delete(publishers, Hash(an.publishValid.solution))
	publishers[an.publishValid.nextPuzzle] = name

	admins.Remove(Hash(an.adminValid.solution))
	admins.Insert(an.adminValid.nextPuzzle)
}
