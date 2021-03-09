package block_chain

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang-collections/collections/set"
)

type Action string

const (
	genesis      = "Genesis"
	post         = "Post"
	changeName   = "ChangeName"
	addPublisher = "AddPublisher"
	addNode      = "AddNode"
)

// These values need to be tweaked to increase security

const HashSize = sha256.Size
const SolutionSize = sha256.Size
const PuzzleSize = sha256.Size

func Hash(a interface{}) [SolutionSize]byte {
	return sha256.Sum256([]byte(fmt.Sprintf("%v", a)))
}

type Block interface {
	GetHash() [HashSize]byte
	GetPrevHash() [HashSize]byte
	SetHash(new_hash [HashSize]byte)
	CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool
	ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set)
}

// base type to be inherited
type block struct {
	prevHash     [HashSize]byte
	hash         [HashSize]byte
	name         string
	publishValid PrivValidation
}

type PrivValidation struct {
	solution   [SolutionSize]byte
	nextPuzzle [PuzzleSize]byte
}

func (b *block) GetHash() [HashSize]byte {
	return b.hash
}

func (b *block) GetPrevHash() [HashSize]byte {
	return b.prevHash
}

func (b *block) SetHash(new_hash [HashSize]byte) {
	b.hash = new_hash
}

func (b *block) CheckValidations(publishers map[[PuzzleSize]byte]string) bool {
	puzzle := Hash(b.publishValid.solution)
	if _, ok := publishers[puzzle]; ok == false {
		return false
	} else {
		return true
	}
}

func CheckHash(b Block) bool {
	// make sure hash makes sense
	given_hash := b.GetHash()
	b.SetHash([HashSize]byte{0})
	measured_hash := Hash(b)
	b.SetHash(given_hash)
	if measured_hash != given_hash {
		return false
	}
	return true
}
