package block_chain

import (
	"container/list"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang-collections/collections/set"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Self struct {
	ID   peer.ID
	Name string
}

type BlockChain struct {
	self        Self
	publish_sol [SolutionSize]byte
	admin_sol   [SolutionSize]byte
	chain       list.List
	publishers  map[[PuzzleSize]byte]string
	admins      *set.Set
	nodes       *set.Set
}

func New(self Self) *BlockChain {
	bc := new(BlockChain)
	bc.self = self
	bc.chain.Init()
	bc.admins = set.New()
	bc.nodes = set.New()
	bc.publishers = make(map[[PuzzleSize]byte]string)

	getValidSol(&bc.publish_sol)
	getValidSol(&bc.admin_sol)

	return bc
}

func getValidSol(sol *[SolutionSize]byte) [SolutionSize]byte {
	var ret [SolutionSize]byte
	copy(ret[:], sol[:])
	n, err := rand.Read(sol[:])
	if err != nil {
		panic(err)
	}
	if n != SolutionSize {
		err_msg := fmt.Sprintf("Wrong size solution Generated for Publish Solution. Expected %d, got %d", SolutionSize, n)
		panic(errors.New(err_msg))
	}
	return ret
}

func (bc *BlockChain) AddBlock(b Block) error {
	// make sure valid hash
	ret := CheckHash(b)
	if ret != true {
		err_msg := fmt.Sprintf("Invalid hash: %v", b)
		return errors.New(err_msg)
	}
	// make sure that validations are correct
	ret = b.CheckValidations(bc.publishers, bc.admins)
	if ret != true {
		err_msg := fmt.Sprintf("Invalid validation: %v", b)
		return errors.New(err_msg)
	}
	// given both of these, lets make sure that prev hash is right (unless Genesis block)
	switch b.(type) {
	case *Genesis:
		if bc.chain.Len() != 0 {
			err_msg := fmt.Sprintf("Attempted Genesis block when chain not empty: %v", bc.chain)
			return errors.New(err_msg)
		}
	default:
		if bc.chain.Len() == 0 {
			err_msg := fmt.Sprintf("Chain has no prev blocks: %v", bc.chain.Back())
			return errors.New(err_msg)
		} else if bc.chain.Back().Value.(Block).GetHash() != b.GetPrevHash() {
			err_msg := fmt.Sprintf("Invalid prevHash: %v, %v", bc.chain.Back(), b)
			return errors.New(err_msg)
		}
	}
	bc.chain.PushBack(b)
	b.ApplyValidations(bc.publishers, bc.admins)
	return nil
}

func (bc *BlockChain) Genesis(title string) error {
	var prev_hash [HashSize]byte
	binary.PutUvarint(prev_hash[:], 37) // no signifigance in this number

	var solution_placeholder [SolutionSize]byte
	binary.PutUvarint(solution_placeholder[:], 37) // no signifigance in this number

	var gen Block
	gen = NewGenesis(prev_hash,
		bc.self.Name,
		PrivValidation{solution: solution_placeholder, nextPuzzle: Hash(bc.publish_sol)},
		PrivValidation{solution: solution_placeholder, nextPuzzle: Hash(bc.admin_sol)},
		title,
		bc.self.ID)

	err := bc.AddBlock(gen)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BlockChain) Post(msg string) error {
	sol := getValidSol(&(bc.publish_sol))
	next_puzzle := Hash(bc.publish_sol)
	pv := PrivValidation{solution: sol, nextPuzzle: next_puzzle}
	p := NewPost(bc.chain.Back().Value.(Block).GetHash(), bc.self.Name, pv, msg)
	err := bc.AddBlock(p)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BlockChain) ChangeName(new_name string) error {
	sol := getValidSol(&(bc.publish_sol))
	next_puzzle := Hash(bc.publish_sol)
	pv := PrivValidation{solution: sol, nextPuzzle: next_puzzle}
	var change_name Block
	change_name = NewChangeName(bc.chain.Back().Value.(Block).GetHash(), bc.self.Name, pv, new_name)
	err := bc.AddBlock(change_name)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BlockChain) AddPublisher(friend_puzzle [PuzzleSize]byte, friend_name string) error {
	sol := getValidSol(&(bc.publish_sol))
	next_puzzle := Hash(bc.publish_sol)
	pv := PrivValidation{solution: sol, nextPuzzle: next_puzzle}

	a_sol := getValidSol(&(bc.admin_sol))
	a_next_puzzle := Hash(bc.admin_sol)
	a_pv := PrivValidation{solution: a_sol, nextPuzzle: a_next_puzzle}

	var add_publisher Block
	add_publisher = NewAddPublisher(bc.chain.Back().Value.(Block).GetHash(),
		bc.self.Name,
		pv,
		a_pv,
		friend_puzzle,
		friend_name)

	add_publisher.CheckValidations(bc.publishers, bc.admins)

	err := bc.AddBlock(add_publisher)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BlockChain) AddNode(friend_peerID peer.ID) error {
	var add_node Block
	add_node = NewAddNode(bc.chain.Back().Value.(Block).GetHash(),
		bc.self.Name,
		PrivValidation{solution: getValidSol(&bc.publish_sol), nextPuzzle: Hash(bc.publish_sol)},
		PrivValidation{solution: getValidSol(&bc.admin_sol), nextPuzzle: Hash(bc.admin_sol)},
		friend_peerID)
	err := bc.AddBlock(add_node)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BlockChain) ShareChain() list.List {
	chain := bc.chain
	return chain
}

func (bc *BlockChain) SharePubPuzzle() [PuzzleSize]byte {
	return Hash(bc.publish_sol[:])
}

func (bc *BlockChain) ShareId() peer.ID {
	return bc.self.ID
}
