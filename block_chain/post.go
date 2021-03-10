package block_chain

import (
	"github.com/golang-collections/collections/set"
)

type Post struct {
	block
	action Action
	msg    string
}

func NewPost(prev_hash [HashSize]byte,
	name string,
	pub_valid PrivValidation,
	msg string) *Post {

	post := Post{
		action: post,
		msg:    msg,
	}
	post.prevHash = prev_hash
	post.name = name
	post.publishValid = pub_valid
	post.hash = [HashSize]byte{0}
	post.hash = Hash(&post)
	return &post
}

func (p *Post) GetHash() [HashSize]byte {
	return p.block.GetHash()
}

func (p *Post) GetPrevHash() [HashSize]byte {
	return p.block.GetPrevHash()
}

func (p *Post) SetHash(new_hash [HashSize]byte) {
	p.block.SetHash(new_hash)
}

func (p *Post) CheckValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) bool {
	if p.block.CheckValidations(publishers) == false {
		return false
	}
	if publishers[Hash(p.publishValid.solution)] != p.name {
		return false
	}
	return true
}

func (p *Post) ApplyValidations(publishers map[[PuzzleSize]byte]string, admins *set.Set) {
	name := publishers[Hash(p.publishValid.solution)]
	delete(publishers, Hash(p.publishValid.solution))
	publishers[p.publishValid.nextPuzzle] = name
}
