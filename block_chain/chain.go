package block_chain

import (
	"container/list"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/golang-collections/collections/set"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

type self struct {
	priv crypto.PrivKey
	iD   peer.ID
	name string
}

func MakeSelf(name string) self {
	priv, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		panic(err)
	}
	id, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		panic(err)
	}
	return self{priv: priv, iD: id, name: name}
}

type BlockChain struct {
	self        self
	publish_sol [SolutionSize]byte
	admin_sol   [SolutionSize]byte

	chain     list.List
	chain_in  chan Block
	chain_err chan error

	asyncOut bool
	ChainOut chan list.List

	publishers map[[PuzzleSize]byte]string
	admins     *set.Set
	nodes      *set.Set
	p2p        *p2pStuff
}

func New(self self, asyncOut bool) *BlockChain {
	bc := new(BlockChain)
	bc.self = self
	bc.chain.Init()
	bc.admins = set.New()
	bc.nodes = set.New()
	bc.publishers = make(map[[PuzzleSize]byte]string)

	bc.p2p = nil

	bc.chain_in = make(chan Block)
	bc.chain_err = make(chan error)

	bc.asyncOut = asyncOut
	bc.ChainOut = make(chan list.List, 1) //SyncValueNew(bc.chain)

	getValidSol(&bc.publish_sol)
	getValidSol(&bc.admin_sol)

	go bc.Run()
	return bc
}

func (bc *BlockChain) Invite() (string, error) {
	if bc.chain.Len() == 0 {
		return "", errors.New("Chain has no blocks")
	}
	return fmt.Sprintf("%s %x", bc.chain.Front().Value.(*Genesis).Title,
		bc.chain.Front().Value.(Block).GetHash()), nil
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

func (bc *BlockChain) Run() {
	for {
		b := <-bc.chain_in
		bc.chain_err <- bc.addBlock(b)
	}
}

func (bc *BlockChain) AddBlock(b Block) error {
	bc.chain_in <- b
	err := <-bc.chain_err
	if err != nil {
		return err
	}
	if bc.p2p != nil {
		bc.p2p.publishBlock(b)
	}
	if bc.asyncOut == true {
		bc.ChainOut <- bc.chain
	}
	return nil
}

func (bc *BlockChain) addBlock(b Block) error {
	// make sure valid hash
	ret := CheckHash(b)
	if ret != true {
		err_msg := fmt.Sprintf("Invalid hash: %+v", b)
		return errors.New(err_msg)
	}
	// make sure that validations are correct
	ret = b.CheckValidations(bc.publishers, bc.admins)
	if ret != true {
		err_msg := fmt.Sprintf("Invalid validation: %+v", b)
		return errors.New(err_msg)
	}
	// given both of these, lets make sure that prev hash is right (unless Genesis block)
	switch b.(type) {
	case *Genesis:
		if bc.chain.Len() != 0 {
			err_msg := fmt.Sprintf("Attempted Genesis block when chain not empty: %+v", bc.chain)
			return errors.New(err_msg)
		}
	default:
		if bc.chain.Len() == 0 {
			err_msg := fmt.Sprintf("Chain has no prev blocks: %+v", bc.chain.Back())
			return errors.New(err_msg)
		} else if bc.chain.Back().Value.(Block).GetHash() != b.GetPrevHash() {
			err_msg := fmt.Sprintf("Invalid prevHash: %+v,\n%+v", bc.chain.Back(), b)
			return errors.New(err_msg)
		}
	}
	bc.chain.PushBack(b)
	b.ApplyValidations(bc.publishers, bc.admins)
	switch b.(type) {
	case *ChangeName:
		cn := b.(*ChangeName)
		if bc.self.name == cn.Name {
			bc.self.name = cn.NewName
		}
	}
	return nil
}

func (bc *BlockChain) Genesis(title string) error {
	var prev_hash [HashSize]byte
	binary.PutUvarint(prev_hash[:], 37) // no signifigance in this number

	var solution_placeholder [SolutionSize]byte
	binary.PutUvarint(solution_placeholder[:], 37) // no signifigance in this number

	var gen Block
	gen = NewGenesis(prev_hash,
		bc.self.name,
		PrivValidation{Solution: solution_placeholder, NextPuzzle: Hash(bc.publish_sol)},
		PrivValidation{Solution: solution_placeholder, NextPuzzle: Hash(bc.admin_sol)},
		title,
		bc.self.iD)

	err := bc.AddBlock(gen)
	if err != nil {
		return err
	}
	bc.Join(title, gen.GetHash())
	return nil
}

func (bc *BlockChain) Post(msg string) error {
	sol := getValidSol(&(bc.publish_sol))
	next_puzzle := Hash(bc.publish_sol)
	pv := PrivValidation{Solution: sol, NextPuzzle: next_puzzle}
	p := NewPost(bc.chain.Back().Value.(Block).GetHash(), bc.self.name, pv, msg)
	err := bc.AddBlock(p)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BlockChain) ChangeName(new_name string) error {
	sol := getValidSol(&(bc.publish_sol))
	next_puzzle := Hash(bc.publish_sol)
	pv := PrivValidation{Solution: sol, NextPuzzle: next_puzzle}
	var change_name Block
	change_name = NewChangeName(bc.chain.Back().Value.(Block).GetHash(), bc.self.name, pv, new_name)
	err := bc.AddBlock(change_name)
	if err != nil {
		return err
	}
	return nil
}

func (bc *BlockChain) AddPublisher(friend_puzzle [PuzzleSize]byte, friend_name string) error {
	sol := getValidSol(&(bc.publish_sol))
	next_puzzle := Hash(bc.publish_sol)
	pv := PrivValidation{Solution: sol, NextPuzzle: next_puzzle}

	a_sol := getValidSol(&(bc.admin_sol))
	a_next_puzzle := Hash(bc.admin_sol)
	a_pv := PrivValidation{Solution: a_sol, NextPuzzle: a_next_puzzle}

	var add_publisher Block
	add_publisher = NewAddPublisher(bc.chain.Back().Value.(Block).GetHash(),
		bc.self.name,
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
		bc.self.name,
		PrivValidation{Solution: getValidSol(&bc.publish_sol), NextPuzzle: Hash(bc.publish_sol)},
		PrivValidation{Solution: getValidSol(&bc.admin_sol), NextPuzzle: Hash(bc.admin_sol)},
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

func (bc *BlockChain) ShareChainSince(hash [HashSize]byte) (list.List, error) {
	var ret list.List

	b := bc.chain.Back()
	for ; b != nil; b = b.Prev() {
		if b.Value.(Block).GetHash() != hash {
			ret.PushFront(b.Value.(Block))
		} else {
			break
		}
	}

	if b == nil {
		err_msg := fmt.Sprintf("Hash %x was not found.", hash)
		return ret, errors.New(err_msg)
	}

	return ret, nil
}

func (bc *BlockChain) SharePubPuzzle() [PuzzleSize]byte {
	return Hash(bc.publish_sol[:])
}

func (bc *BlockChain) ShareID() peer.ID {
	return bc.self.iD
}
