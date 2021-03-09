package block_chain

import (
	"fmt"
	"github.com/golang-collections/collections/set"
)

type ChangeName struct {
	block
	NewName string
}

func NewChangeName(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	new_name string) *ChangeName {

	cn := ChangeName{
		NewName: new_name,
	}
	cn.PrevHash = prev_hash
	cn.Name = name
	cn.PublishValid = pub_valid
	cn.Action = changeName
	cn.Hash = [HashSize]byte{0}
	cn.Hash = Hash(&cn)
	return &cn
}

func (cn *ChangeName) GetAction() Action {
	return cn.block.GetAction()
}

func (cn *ChangeName) AsString() string {
	return fmt.Sprintf("%s Changed Name to %s\n", cn.Name, cn.NewName)
}

func (cn *ChangeName) GetHash() [HashSize]byte {
	return cn.block.GetHash()
}

func (cn *ChangeName) GetPrevHash() [HashSize]byte {
	return cn.block.GetPrevHash()
}

func (cn *ChangeName) GetName() string {
	return cn.block.GetName()
}

func (cn *ChangeName) SetHash(new_hash [HashSize]byte) {
	cn.block.SetHash(new_hash)
}

func (cn *ChangeName) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	if cn.block.CheckValidations(publishers) == false {
		return false
	}
	if publishers[Hash(cn.PublishValid.Solution)] != cn.Name {
		return false
	}
	return true
}

func (cn *ChangeName) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	delete(publishers, Hash(cn.PublishValid.Solution))
	publishers[cn.PublishValid.NextPuzzle] = cn.NewName
}
