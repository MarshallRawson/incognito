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
	Id   peer.ID
	Name string
}

type BlockChain struct {
	self      Self
	pub_sol   [SolutionSize]byte
	priv_sol  [SolutionSize]byte
	chain     list.List                   // the actual list of blocks
	pub_name  map[[PuzzleSize]byte]string // people with solutions to this are able to post messages
	priv      *set.Set                    // people with solutions to this AND a pub are able to add/remove pubs/vers
	verifiers *set.Set                    // these are peers who are sent every block in the chain and are asked to verify
}

func New(self Self) *BlockChain {
	bc := new(BlockChain)
	bc.self = self
	bc.chain.Init()
	bc.priv = set.New()
	bc.verifiers = set.New()
	bc.pub_name = make(map[[PuzzleSize]byte]string)

	n, err := rand.Read(bc.pub_sol[:])
	if err != nil {
		panic(err)
	}
	if n != SolutionSize {
		err_msg := fmt.Sprintf("Wrong size solution Generated for Publish Solution. Expected %d, got %d", SolutionSize, n)
		panic(errors.New(err_msg))
	}

	n, err = rand.Read(bc.priv_sol[:])
	if err != nil {
		panic(err)
	}
	if n != SolutionSize {
		err_msg := fmt.Sprintf("Wrong size solution Generated for Privileged Solution. Expected %d, got %d", SolutionSize, n)
		panic(errors.New(err_msg))
	}
	return bc
}

func (bc *BlockChain) Genesis(title string) {
	var hash [HashSize]byte
	binary.PutUvarint(hash[:], 37) // no signifigance in this number

	var pub_sol [SolutionSize]byte
	binary.PutUvarint(pub_sol[:], 37) // no signifigance in this number

	block := Block{
		PrevHash: hash,
		PubSol:   pub_sol,
		Name:     bc.self.Name,
		NextPub:  Hash(bc.pub_sol[:]),
	}
	block.MakeGenesis(title, Hash(bc.priv_sol[:]), bc.self.Id)
	err := bc.AddBlock(block)
	if err != nil {
		panic(err)
	}
}

func (bc *BlockChain) AddBlock(block Block) error {
	ret, err := bc.CheckBlock(block)
	if err != nil {
		return err
	}
	if ret != Success {
		err_msg := fmt.Sprintf("Invalid block: %s", ret)
		return errors.New(err_msg)
	}

	switch block.Action {
	case Genesis:
		bc.addGenesisBlock(block)
	case Post:
		bc.addPostBlock(block)
	case NameChange:
		bc.addNameChangeBlock(block)
	case AddPub:
		bc.addAddPubBlock(block)
	case AddVer:
		bc.addAddVerBlock(block)
	default:
		err_msg := fmt.Sprintf("Invalid block action: %s", block.Action)
		return errors.New(err_msg)
	}
	return nil
}

type BlockCheckReturn string

const (
	Success          = "Success"          // No Failure, the block in question fits!
	NoHistory        = "NoHistory"        // Attempting to put non-genesis block as not-first
	BadName          = "BadName"          // Attempting to put non-genesis block as not-first
	GenesisMisplaced = "GenesisMisplaced" // Attempting to put genesis block as not-first
	BadPrevHash      = "BadPrevHash"      // Prev hash does not match the hash of the prev message
	BadHash          = "BadHash"          // The hash of this block is incorrect
	BadPubSol        = "BadPubSol"        // The Publisher Solution for this block is wrong
	BadPrivSol       = "BadPrivSol"       // The Privleged Solution for this block is wrong
	Misc             = "Misc"             // Some other error occured (return non-nil error)
)

func (bc *BlockChain) CheckBlock(block Block) (BlockCheckReturn, error) {
	if block.Action == Genesis {
		if bc.chain.Len() != 0 {
			return GenesisMisplaced, nil
		}
		if block.Hash != Hash(block.serialize()) {
			return BadHash, nil
		}
	} else if block.Action == Post ||
		block.Action == NameChange ||
		block.Action == AddPub ||
		block.Action == AddVer {
		if bc.chain.Len() == 0 {
			return NoHistory, nil
		}

		// make sure fits in histroy
		if block.PrevHash != bc.chain.Back().Value.(Block).Hash {
			return BadPrevHash, nil
		}

		// make sure has valid hash
		if block.Hash != Hash(block.serialize()) {
			return BadHash, nil
		}

		// proof of authorized publisher
		if _, ok := bc.pub_name[Hash(block.PubSol[:])]; ok == false {
			return BadPubSol, nil
		}

		// make sure the pub name is consistent
		if name, _ := bc.pub_name[Hash(block.PubSol[:])]; name != block.Name {
			return BadName, nil
		}

		// if adding a publisher, we need to check the privileged solution as well
		if block.Action == AddPub {
			add_pub, err := block.FromAddPub()
			if err != nil {
				panic(err)
			}
			if bc.priv.Has(Hash(add_pub.PrivSol[:])) == false {
				return BadPrivSol, nil
			}
		}
		// if adding a verifier, we also need to check the privileged
		if block.Action == AddVer {
			add_ver, err := block.FromAddVer()
			if err != nil {
				panic(err)
			}
			if bc.priv.Has(Hash(add_ver.PrivSol[:])) == false {
				return BadPrivSol, nil
			}
		}
	} else {
		err_msg := fmt.Sprintf("Unrecognized block type %T!\n", block)
		return Misc, errors.New(err_msg)
	}
	return Success, nil
}

func (bc *BlockChain) addGenesisBlock(block Block) {
	// does not check hashes on block chain
	gen, err := block.FromGenesis()
	if err != nil {
		panic(err)
	}
	bc.chain.PushBack(block)
	bc.verifiers.Insert(gen.Ver)
	bc.pub_name[block.NextPub] = block.Name
	bc.priv.Insert(gen.Priv)
}

func (bc *BlockChain) updateName(block Block) {
	name := bc.pub_name[Hash(block.PubSol[:])]
	delete(bc.pub_name, Hash(block.PubSol[:]))
	bc.pub_name[block.NextPub] = name
}

func (bc *BlockChain) addPostBlock(block Block) {
	// does not check hashes on block chain
	_, err := block.FromPost()
	if err != nil {
		panic(err)
	}
	bc.chain.PushBack(block)
	bc.updateName(block)
}

func (bc *BlockChain) addNameChangeBlock(block Block) {
	// does not check hashes on block chain
	new_name, err := block.FromNameChange()
	if err != nil {
		panic(err)
	}
	bc.chain.PushBack(block)
	bc.updateName(block)
	bc.pub_name[block.NextPub] = new_name
}

func (bc *BlockChain) addAddPubBlock(block Block) {
	// does not check hashes on block chain
	ap, err := block.FromAddPub()
	if err != nil {
		panic(err)
	}
	bc.chain.PushBack(block)
	bc.updateName(block)
	bc.pub_name[ap.NewPub] = ap.NewName
	bc.priv.Remove(Hash(ap.PrivSol[:]))
	bc.priv.Insert(ap.NextPriv)
}

func (bc *BlockChain) addAddVerBlock(block Block) {
	// does not check hashes on block chain
	av, err := block.FromAddVer()
	if err != nil {
		panic(err)
	}
	bc.chain.PushBack(block)
	bc.updateName(block)
	bc.priv.Remove(Hash(av.PrivSol[:]))
	bc.priv.Insert(av.NextPriv)
	bc.verifiers.Insert(av.NewVer)
}

func (bc *BlockChain) MakeNextBlock() Block {
	block := Block{Name: bc.self.Name}
	block.PrevHash = bc.chain.Back().Value.(Block).Hash
	block.PubSol = bc.pub_sol
	_, err := rand.Read(bc.pub_sol[:]) // replenish the pub_sol
	if err != nil {
		panic(err)
	}
	block.NextPub = Hash(bc.pub_sol[:])
	return block
}

func (bc *BlockChain) Post(msg string) {
	post := bc.MakeNextBlock()
	post.MakePost(msg)
	err := bc.AddBlock(post)
	if err != nil {
		panic(err)
	}
}

func (bc *BlockChain) NameChange(new_name string) {
	name_change := bc.MakeNextBlock()
	name_change.MakeNameChange(new_name)
	err := bc.AddBlock(name_change)
	if err != nil {
		panic(err)
	}
}

func (bc *BlockChain) getPriv() ([SolutionSize]byte, [PuzzleSize]byte) {
	priv_sol := bc.priv_sol

	_, err := rand.Read(bc.priv_sol[:])
	if err != nil {
		panic(err)
	}
	return priv_sol, Hash(bc.priv_sol[:])
}

func (bc *BlockChain) AddPub(friend_puzzle [PuzzleSize]byte, friend_name string) {
	friend := bc.MakeNextBlock()
	priv_sol, next_priv := bc.getPriv()
	friend.MakeAddPub(AddPubExtras{friend_puzzle, friend_name, priv_sol, next_priv})
	err := bc.AddBlock(friend)
	if err != nil {
		panic(err)
	}
}

func (bc *BlockChain) AddVer(friend_peerID peer.ID) {
	friend := bc.MakeNextBlock()
	priv_sol, next_priv := bc.getPriv()
	friend.MakeAddVer(AddVerExtras{priv_sol, next_priv, friend_peerID})
	err := bc.AddBlock(friend)
	if err != nil {
		panic(err)
	}
}

func (bc *BlockChain) ShareChain() list.List {
	chain := bc.chain
	return chain
}

func (bc *BlockChain) SharePubPuzzle() [PuzzleSize]byte {
	return Hash(bc.pub_sol[:])
}

func (bc *BlockChain) ShareId() peer.ID {
	return bc.self.Id
}
