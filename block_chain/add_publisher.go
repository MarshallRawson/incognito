package block_chain

import (
	"fmt"
	"github.com/golang-collections/collections/set"
)

type AddPublisher struct {
	block
	AdminValid         PrivValidation
	NewPublisherPuzzle [PuzzleSize]byte
	NewName            string
}

func NewAddPublisher(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	admin_valid PrivValidation,
	new_publisher_puzzle [PuzzleSize]byte,
	new_name string) *AddPublisher {

	ap := AddPublisher{
		AdminValid:         admin_valid,
		NewPublisherPuzzle: new_publisher_puzzle,
		NewName:            new_name,
	}
	ap.PrevHash = prev_hash
	ap.Name = name
	ap.PublishValid = pub_valid
	ap.Action = addPublisher
	ap.Hash = [HashSize]byte{0}
	ap.Hash = Hash(&ap)
	return &ap
}

func (ap *AddPublisher) GetAction() Action {
	return ap.block.GetAction()
}

func (ap *AddPublisher) AsString() string {
	return fmt.Sprintf("%s: Added Publisher: %s\n", ap.Name, ap.NewName)
}

func (ap *AddPublisher) GetHash() [HashSize]byte {
	return ap.block.GetHash()
}

func (ap *AddPublisher) GetPrevHash() [HashSize]byte {
	return ap.block.GetPrevHash()
}

func (ap *AddPublisher) GetName() string {
	return ap.block.GetName()
}

func (ap *AddPublisher) SetHash(new_hash [HashSize]byte) {
	ap.block.SetHash(new_hash)
}

func (ap *AddPublisher) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	if ap.block.CheckValidations(publishers) == false {
		return false
	}
	if publishers[Hash(ap.PublishValid.Solution)] != ap.Name {
		return false
	}
	if admins.Has(Hash(ap.AdminValid.Solution)) == false {
		return false
	}
	return true
}

func (ap *AddPublisher) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	name := publishers[Hash(ap.PublishValid.Solution)]
	delete(publishers, Hash(ap.PublishValid.Solution))
	publishers[ap.PublishValid.NextPuzzle] = name

	publishers[ap.NewPublisherPuzzle] = ap.NewName

	admins.Remove(Hash(ap.AdminValid.Solution))
	admins.Insert(ap.AdminValid.NextPuzzle)
}
