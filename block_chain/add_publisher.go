package block_chain

import (
	"github.com/golang-collections/collections/set"
)

type AddPublisher struct {
	block
	action             Action
	adminValid         PrivValidation
	newPublisherPuzzle [PuzzleSize]byte
	newName            string
}

func MakeAddPublisher(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	admin_valid PrivValidation,
	new_publisher_puzzle [PuzzleSize]byte,
	new_name string) *AddPublisher {

	ap := AddPublisher{
		action:             addPublisher,
		adminValid:         admin_valid,
		newPublisherPuzzle: new_publisher_puzzle,
		newName:            new_name,
	}
	ap.prevHash = prev_hash
	ap.name = name
	ap.publishValid = pub_valid
	ap.hash = [HashSize]byte{0}
	ap.hash = Hash(&ap)
	return &ap
}

func (ap *AddPublisher) GetHash() [HashSize]byte {
	return ap.block.GetHash()
}

func (ap *AddPublisher) GetPrevHash() [HashSize]byte {
	return ap.block.GetPrevHash()
}

func (ap *AddPublisher) SetHash(new_hash [HashSize]byte) {
	ap.block.SetHash(new_hash)
}

func (ap *AddPublisher) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	if ap.block.CheckValidations(publishers) == false {
		return false
	}
	if publishers[Hash(ap.publishValid.solution)] != ap.name {
		return false
	}
	if admins.Has(Hash(ap.adminValid.solution)) == false {
		return false
	}
	return true
}

func (ap *AddPublisher) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	name := publishers[Hash(ap.publishValid.solution)]
	delete(publishers, Hash(ap.publishValid.solution))
	publishers[ap.publishValid.nextPuzzle] = name

	publishers[ap.newPublisherPuzzle] = ap.newName

	admins.Remove(Hash(ap.adminValid.solution))
	admins.Insert(ap.adminValid.nextPuzzle)
}
