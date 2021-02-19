package block_chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
)

type BlockAction string

const (
	Genesis    = "Genesis"
	Post       = "Post"
	NameChange = "NameChange"
	AddPub     = "AddPub"
	AddVer     = "AddVer"
	AddPriv    = "AddPriv" //TODO (if we want this option)
)

// These values need to be tweaked to increase security

const HashSize = sha256.Size
const SolutionSize = sha256.Size
const PuzzleSize = sha256.Size

func Hash(a []byte) [SolutionSize]byte {
	return sha256.Sum256(a)
}

// base type to be inherited
type Block struct {
	Hash [HashSize]byte // sha256.Sum256 of all the other elements

	Action   BlockAction        // Elements in all Blocks
	PrevHash [HashSize]byte     // sha256.Sum256 of the previous message (excluding the previous messages hash)
	PubSol   [SolutionSize]byte // solution to a public to prove authorized
	NextPub  [PuzzleSize]byte   // next pub hash puzzel that proves identity
	Name     string
	Msg      []byte
}

func (block Block) AsStringVerbose() string {
	ret := fmt.Sprintf("{\n")
	ret += fmt.Sprintf("  Previous Hash: %x\n", block.PrevHash)
	ret += fmt.Sprintf("  Publisher Verification: %x", block.PubSol)
	ret += fmt.Sprintf("  Next Publisher Verification: %x", block.PubSol)
	ret += fmt.Sprintf("  Author Name: %s\n", block.Name)
	ret += fmt.Sprintf("  Action: %s\n", block.Action)

	var ex interface{}
	switch block.Action {
	case Genesis:
		ex, _ = block.FromGenesis()
	case Post:
		ex, _ = block.FromPost()
	case NameChange:
		ex, _ = block.FromNameChange()
	case AddPub:
		ex, _ = block.FromAddPub()
	case AddVer:
		ex, _ = block.FromAddVer()
	default:
		err_msg := fmt.Sprintf("Invalid block action: %s", block.Action)
		panic(errors.New(err_msg))
	}
	ex_str, err := json.MarshalIndent(ex, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	ret += string(ex_str)
	ret += fmt.Sprintf("  Hash: %x\n", block.Hash)
	ret += fmt.Sprintf("}")
	return ret
}

func (b *Block) serialize() []byte {
	var arr []byte
	arr = append(arr, []byte(b.Action)...)
	arr = append(arr, 0)
	arr = append(arr, b.PrevHash[:]...)
	arr = append(arr, b.PubSol[:]...)
	arr = append(arr, b.NextPub[:]...)
	arr = append(arr, b.Name...)
	arr = append(arr, b.Msg...)
	return arr
}

func (b *Block) MakePost(msg string) {
	b.Action = Post
	b.Msg = []byte(msg)
	b.Hash = Hash(b.serialize())
}

func (b *Block) FromPost() (string, error) {
	if b.Action != Post {
		return "", errors.New("Block Action is not Post")
	}
	return string(b.Msg), nil
}

func (b *Block) MakeNameChange(new_name string) {
	b.Action = NameChange
	b.Msg = []byte(new_name)
	b.Hash = Hash(b.serialize())
}

func (b *Block) FromNameChange() (string, error) {
	if b.Action != NameChange {
		return "", errors.New("Block Action is not NameChange")
	}
	return string(b.Msg), nil
}

type GenesisExtras struct {
	Title string
	Priv  [PuzzleSize]byte
	Ver   peer.ID
}

func (b *Block) MakeGenesis(title string, priv [PuzzleSize]byte, ver peer.ID) {
	b.Action = Genesis
	b.Msg = b.Msg[:0]
	b.Msg = append(b.Msg, []byte(title)...)
	b.Msg = append(b.Msg, 0)
	b.Msg = append(b.Msg, priv[:]...)
	b.Msg = append(b.Msg, ver[:]...)
	b.Hash = Hash(b.serialize())
}

func (b *Block) FromGenesis() (GenesisExtras, error) {
	gen := GenesisExtras{}
	if b.Action != Genesis {
		return gen, errors.New("Block Action is not Genesis")
	}
	seek := bytes.Index(b.Msg, []byte{0})
	gen.Title = string(b.Msg[:seek])
	seek = seek + 1
	for i, val := range b.Msg[seek : seek+PuzzleSize] {
		gen.Priv[i] = val
	}
	seek = seek + PuzzleSize
	gen.Ver = peer.ID(b.Msg[seek:])
	return gen, nil
}

type AddPubExtras struct {
	NewPub   [PuzzleSize]byte
	NewName  string
	PrivSol  [SolutionSize]byte
	NextPriv [PuzzleSize]byte
}

func (b *Block) MakeAddPub(add_pub AddPubExtras) {
	b.Action = AddPub
	b.Msg = b.Msg[:0]
	b.Msg = append(b.Msg, add_pub.NewPub[:]...)
	b.Msg = append(b.Msg, []byte(add_pub.NewName)...)
	b.Msg = append(b.Msg, 0)
	b.Msg = append(b.Msg, add_pub.PrivSol[:]...)
	b.Msg = append(b.Msg, add_pub.NextPriv[:]...)
	b.Hash = Hash(b.serialize())
}

func (b *Block) FromAddPub() (AddPubExtras, error) {
	ap := AddPubExtras{}
	if b.Action != AddPub {
		return ap, errors.New("Block Action is not AddPub")
	}
	seek := PuzzleSize
	for i, val := range b.Msg[:seek] {
		ap.NewPub[i] = val
	}
	null_term := bytes.Index(b.Msg[seek:], []byte{0}) + seek
	ap.NewName = string(b.Msg[seek:null_term])
	seek = null_term + 1
	for i, val := range b.Msg[seek : seek+SolutionSize] {
		ap.PrivSol[i] = val
	}
	seek = seek + PuzzleSize
	for i, val := range b.Msg[seek:] {
		ap.NextPriv[i] = val
	}
	return ap, nil
}

type AddVerExtras struct {
	PrivSol  [SolutionSize]byte
	NextPriv [PuzzleSize]byte
	NewVer   peer.ID
}

func (b *Block) MakeAddVer(add_ver AddVerExtras) {
	b.Action = AddVer
	b.Msg = b.Msg[:0]
	b.Msg = append(b.Msg, add_ver.PrivSol[:]...)
	b.Msg = append(b.Msg, add_ver.NextPriv[:]...)
	b.Msg = append(b.Msg, add_ver.NewVer[:]...)
	b.Hash = Hash(b.serialize())
}

func (b *Block) FromAddVer() (AddVerExtras, error) {
	av := AddVerExtras{}
	if b.Action != AddVer {
		return av, errors.New("Block Action is not AddVer")
	}
	seek := 0
	for i, val := range b.Msg[seek : seek+SolutionSize] {
		av.PrivSol[i] = val
	}
	seek = seek + SolutionSize
	for i, val := range b.Msg[seek : seek+PuzzleSize] {
		av.NextPriv[i] = val
	}
	seek = seek + PuzzleSize
	av.NewVer = peer.ID(b.Msg[seek:])
	return av, nil
}
