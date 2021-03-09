package block_chain

import (
	"crypto/sha256"
	"encoding/json"
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
	GetAction() Action
	GetName() string
	GetPrevHash() [HashSize]byte
	SetHash(new_hash [HashSize]byte)
	CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool
	ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set)
	AsString() string
}

// base type to be inherited
type block struct {
	PrevHash     [HashSize]byte
	Hash         [HashSize]byte
	Name         string
	PublishValid PrivValidation
	Action       Action
}

type PrivValidation struct {
	Solution   [SolutionSize]byte
	NextPuzzle [PuzzleSize]byte
}

func (b *block) GetAction() Action {
	return b.Action
}

func (b *block) GetHash() [HashSize]byte {
	return b.Hash
}

func (b *block) GetPrevHash() [HashSize]byte {
	return b.PrevHash
}

func (b *block) GetName() string {
	return b.Name
}

func (b *block) SetHash(new_hash [HashSize]byte) {
	b.Hash = new_hash
}

func (b *block) CheckValidations(publishers map[[PuzzleSize]byte]string) bool {
	puzzle := Hash(b.PublishValid.Solution)
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

func MarshalBlock(b Block) ([]byte, error) {
	var v []byte
	var err error
	switch b.GetAction() {
	case genesis:
		v, err = json.Marshal(*b.(*Genesis))
	case post:
		v, err = json.Marshal(*b.(*Post))
	case changeName:
		v, err = json.Marshal(*b.(*ChangeName))
	case addPublisher:
		v, err = json.Marshal(*b.(*AddPublisher))
	case addNode:
		v, err = json.Marshal(*b.(*AddNode))
	}
	if err != nil {
		return []byte{}, err
	}
	type_label := []byte(string(b.GetAction()))
	d := [][]byte{type_label, v}
	payload, err := json.Marshal(d)
	if err != nil {
		return []byte{}, err
	}
	return payload, nil
}

func UnmarshalBlock(b []byte) (Block, error) {
	d := make([][]byte, 2)
	err := json.Unmarshal(b, &d)
	var ret Block
	if err != nil {
		return ret, err
	}
	switch Action(d[0]) {
	case genesis:
		g := new(Genesis)
		err = json.Unmarshal(d[1], g)
		ret = g
	case post:
		g := new(Post)
		err = json.Unmarshal(d[1], g)
		ret = g
	case changeName:
		g := new(ChangeName)
		err = json.Unmarshal(d[1], g)
		ret = g
	case addPublisher:
		g := new(AddPublisher)
		err = json.Unmarshal(d[1], g)
		ret = g
	case addNode:
		g := new(AddNode)
		err = json.Unmarshal(d[1], g)
		ret = g
	}

	return ret, err
}
