package block_chain

import (
	"fmt"
	"github.com/golang-collections/collections/set"
)

type Post struct {
	block
	Msg string
}

func NewPost(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	msg string) *Post {

	p := Post{
		Msg: msg,
	}
	p.PrevHash = prev_hash
	p.Name = name
	p.PublishValid = pub_valid
	p.Action = post
	p.Hash = [HashSize]byte{0}
	p.Hash = Hash(&p)
	return &p
}

func (p *Post) GetAction() Action {
	return p.block.GetAction()
}

func (p *Post) AsString() string {
	return fmt.Sprintf("%s: %s\n", p.Name, p.Msg)
}

func (p *Post) GetHash() [HashSize]byte {
	return p.block.GetHash()
}

func (p *Post) GetPrevHash() [HashSize]byte {
	return p.block.GetPrevHash()
}

func (p *Post) GetName() string {
	return p.block.GetName()
}

func (p *Post) SetHash(new_hash [HashSize]byte) {
	p.block.SetHash(new_hash)
}

func (p *Post) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	if p.block.CheckValidations(publishers) == false {
		return false
	}
	if publishers[Hash(p.PublishValid.Solution)] != p.Name {
		return false
	}
	return true
}

func (p *Post) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	name := publishers[Hash(p.PublishValid.Solution)]
	delete(publishers, Hash(p.PublishValid.Solution))
	publishers[p.PublishValid.NextPuzzle] = name
}
