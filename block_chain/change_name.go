package block_chain

import (
	"github.com/golang-collections/collections/set"
)

type ChangeName struct {
	block
	action  Action
	newName string
}

func NewChangeName(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	new_name string) *ChangeName {

	cn := ChangeName{
		action:  changeName,
		newName: new_name,
	}
	cn.prevHash = prev_hash
	cn.name = name
	cn.publishValid = pub_valid
	cn.hash = [HashSize]byte{0}
	cn.hash = Hash(&cn)
	return &cn
}

func (cn *ChangeName) GetHash() [HashSize]byte {
	return cn.block.GetHash()
}

func (cn *ChangeName) GetPrevHash() [HashSize]byte {
	return cn.block.GetPrevHash()
}

func (cn *ChangeName) SetHash(new_hash [HashSize]byte) {
	cn.block.SetHash(new_hash)
}

func (cn *ChangeName) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	if cn.block.CheckValidations(publishers) == false {
		return false
	}
	if publishers[Hash(cn.publishValid.solution)] != cn.name {
		return false
	}
	return true
}

func (cn *ChangeName) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	delete(publishers, Hash(cn.publishValid.solution))
	publishers[cn.publishValid.nextPuzzle] = cn.newName
}
