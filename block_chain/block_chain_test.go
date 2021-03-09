package block_chain

import (
	"crypto/rand"
	"github.com/libp2p/go-libp2p-core/peer"
	"testing"
)

// Testing the block chain
func TestBlockChain(t *testing.T) {
	// Alice constructs her block chain
	bc := New(MakeSelf("alice"))
	// Alice makes her own block chain
	err := bc.Genesis("Alice's page")
	if err != nil {
		panic(err)
	}

	// Alice makes a post on her block chain
	err = bc.Post("Hello World!")
	if err != nil {
		panic(err)
	}

	// Alice meets a person named bob
	bob_chain := New(MakeSelf("bob"))

	// Alice acquires a publisher puzzle from Bob, which bob knows how to solve, but Alice does not
	bob_puzzle := bob_chain.SharePubPuzzle()

	// Alice gets bob's public puzzle and makes a bock which adds bob as a publisher
	err = bc.AddPublisher(bob_puzzle, "bob")
	if err != nil {
		panic(err)
	}

	// Alice makes a post
	err = bc.Post("I added Bob")
	if err != nil {
		panic(err)
	}

	// Bob Shares with Alice his peerID, and Alice adds him as a verifier
	// so he will recive updates
	err = bc.AddNode(bob_chain.ShareID())
	if err != nil {
		panic(err)
	}

	// Alice changes her name from "alice" to "Alice"
	err = bc.ChangeName("Alice")
	if err != nil {
		panic(err)
	}

	// check the blocks are valid
	chain := bc.ShareChain()
	b := chain.Front()

	gen := b.Value.(*Genesis)
	if gen.Action != genesis {
		t.Errorf("Expected Action %s, but got %s", genesis, gen.Action)
	}
	if gen.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", gen.Name)
	}
	if gen.Title != "Alice's page" {
		t.Errorf("Expected Title %s, but got %s", "Alice's page", gen.Name)
	}
	hash := gen.Hash
	b = b.Next()

	p := b.Value.(*Post)
	if p.Action != post {
		t.Errorf("Expected Action %s, but got %s", post, p.Action)
	}
	if p.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", p.Name)
	}
	if p.Msg != "Hello World!" {
		t.Errorf("Expected  %s, but got %s", "Hello World!", p.Msg)
	}
	if p.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, p.PrevHash)
	}
	hash = p.Hash
	b = b.Next()

	pub := b.Value.(*AddPublisher)
	if pub.Action != addPublisher {
		t.Errorf("Expected Action %s, but got %s", "AddPub", pub.Action)
	}
	if pub.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", pub.Name)
	}
	if pub.NewName != "bob" {
		t.Errorf("Expected NewName %s, but got %s", "bob", pub.NewName)
	}
	if pub.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, pub.PrevHash)
	}
	if pub.NewPublisherPuzzle != bob_chain.SharePubPuzzle() {
		t.Errorf("Expected  pub puzzle %x, but got %x",
			bob_chain.SharePubPuzzle(), pub.NewPublisherPuzzle)
	}
	hash = pub.Hash
	b = b.Next()

	post2 := b.Value.(*Post)
	if post2.Action != post {
		t.Errorf("Expected Action %s, but got %s", post, post2.Action)
	}
	if post2.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", post2.Name)
	}
	if post2.Msg != "I added Bob" {
		t.Errorf("Expected  %s, but got %s", "I added Bob", post2.Msg)
	}
	if post2.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, post2.PrevHash)
	}
	hash = post2.Hash
	b = b.Next()

	av := b.Value.(*AddNode)
	if av.Action != addNode {
		t.Errorf("Expected Action %s, but got %s", addNode, av.Action)
	}
	if av.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", av.Name)
	}
	if av.NewNodeID != bob_chain.ShareID() {
		t.Errorf("Expected NewVer %x, but got %x", bob_chain.ShareID(), av.NewNodeID)
	}
	if av.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, av.PrevHash)
	}
	hash = av.Hash
	b = b.Next()

	nc := b.Value.(*ChangeName)
	if nc.Action != changeName {
		t.Errorf("Expected Action %s, but got %s", "NameChange", nc.Action)
	}
	if nc.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", nc.Name)
	}
	if nc.NewName != "Alice" {
		t.Errorf("Expected New Name %s, but got %s", "Alice", nc.NewName)
	}
	if nc.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, nc.PrevHash)
	}
	hash = nc.Hash
	b = b.Next()

	// TODO
	// serialize the Whole Block Chain
	// send the chain to bob
	// have bob reconstruct his instance of the chain
	// have bob add a post
	// serialize the block and send it to alice
}

// Testing the Block Methods

func TestGenesisBlock(t *testing.T) {
	name := "Spaceduck"
	title := "Intellectual Dark Web"
	var prev_hash [HashSize]byte
	rand.Read(prev_hash[:])
	pub_valid := PrivValidation{}
	rand.Read(pub_valid.Solution[:])
	rand.Read(pub_valid.NextPuzzle[:])
	admin_valid := PrivValidation{}
	rand.Read(admin_valid.Solution[:])
	rand.Read(admin_valid.NextPuzzle[:])

	gen := NewGenesis(prev_hash, name, pub_valid, admin_valid, title, peer.ID([]byte{3, 7}))

	if gen.Action != genesis {
		t.Errorf("NewGenesis made block action %s, when it should have been %s",
			gen.Action, genesis)
	}

	given_hash := gen.GetHash()
	gen.SetHash([HashSize]byte{0})
	measured_hash := Hash(gen)
	gen.SetHash(given_hash)
	if gen.Hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, gen.Hash)
	}

	idw := string("Intellectual Dark Web")
	if gen.Title != idw {
		t.Errorf("Expected Title %s, but got %s", idw, gen.Title)
	}
	node := peer.ID([]byte{3, 7})
	if gen.NodeID != node {
		t.Errorf("Expected ID %x, but got %x", node, gen.NodeID)
	}
}

func TestPostBlock(t *testing.T) {
	name := "SpaceDuck"
	var prevHash [HashSize]byte
	rand.Read(prevHash[:])
	var pub_valid PrivValidation
	rand.Read(pub_valid.Solution[:])
	rand.Read(pub_valid.NextPuzzle[:])
	p := NewPost(prevHash, name, pub_valid, "Welcome to the Intellecutial Dark Web!")

	if p.Action != post {
		t.Errorf("NewPost made block action %s, when it should have been %s", p.Action, post)
	}

	given_hash := p.GetHash()
	p.SetHash([HashSize]byte{0})
	measured_hash := Hash(p)
	p.SetHash(given_hash)
	if p.Hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, p.Hash)
	}

	idw := string("Welcome to the Intellecutial Dark Web!")
	if p.Msg != idw {
		t.Errorf("Expected Post %s, but got %s", idw, p.Msg)
	}
}

func TestNameChangeBlock(t *testing.T) {
	name := "SpaceDuck"
	var prevHash [HashSize]byte
	rand.Read(prevHash[:])
	var pub_valid PrivValidation
	rand.Read(pub_valid.Solution[:])
	rand.Read(pub_valid.NextPuzzle[:])
	change_name := NewChangeName(prevHash, name, pub_valid, "$paceDuck")

	if change_name.Action != changeName {
		t.Errorf("NewNameChange made block action %s, when it should have been %s",
			change_name.Action, changeName)
	}

	given_hash := change_name.GetHash()
	change_name.SetHash([HashSize]byte{0})
	measured_hash := Hash(change_name)
	change_name.SetHash(given_hash)
	if change_name.Hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, change_name.Hash)
	}

	name = "$paceDuck"
	if change_name.NewName != name {
		t.Errorf("Expected Post %s, but got %s", name, change_name.NewName)
	}
}

func TestAddPubBlock(t *testing.T) {
	name := "SpaceDuck"
	var prev_hash [HashSize]byte
	rand.Read(prev_hash[:])
	pub_valid := PrivValidation{}
	rand.Read(pub_valid.Solution[:])
	rand.Read(pub_valid.NextPuzzle[:])
	admin_valid := PrivValidation{}
	rand.Read(admin_valid.Solution[:])
	rand.Read(admin_valid.NextPuzzle[:])

	var new_pub [PuzzleSize]byte
	rand.Read(new_pub[:])

	ap := NewAddPublisher(prev_hash, name, pub_valid, admin_valid, new_pub, "SpaceGoose")

	if ap.Action != addPublisher {
		t.Errorf("NewPost made block action %s, when it should have been %s", ap.Action, addPublisher)
	}

	given_hash := ap.GetHash()
	ap.SetHash([HashSize]byte{0})
	measured_hash := Hash(ap)
	ap.SetHash(given_hash)
	if ap.Hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, ap.Hash)
	}

	if ap.NewPublisherPuzzle != new_pub {
		t.Errorf("Expected NewPub %x, but got %x", new_pub, ap.NewPublisherPuzzle)
	}
	me := "SpaceGoose"
	if ap.NewName != me {
		t.Errorf("Expected NewName %s, but got %s", me, ap.NewName)
	}
	if ap.AdminValid.Solution != admin_valid.Solution {
		t.Errorf("Expected PrivSol %x, but got %x", admin_valid.Solution, ap.AdminValid.Solution)
	}
	if ap.AdminValid.NextPuzzle != admin_valid.NextPuzzle {
		t.Errorf("Expected PrivSol %x, but got %x", ap.AdminValid.NextPuzzle, admin_valid.NextPuzzle)
	}
}

func TestAddNodeBlock(t *testing.T) {
	name := "SpaceDuck"
	var prev_hash [HashSize]byte
	rand.Read(prev_hash[:])
	pub_valid := PrivValidation{}
	rand.Read(pub_valid.Solution[:])
	rand.Read(pub_valid.NextPuzzle[:])
	admin_valid := PrivValidation{}
	rand.Read(admin_valid.Solution[:])
	rand.Read(admin_valid.NextPuzzle[:])
	node := peer.ID([]byte{3, 7})

	av := NewAddNode(prev_hash, name, pub_valid, admin_valid, node)

	if av.Action != addNode {
		t.Errorf("NewPost made block action %s, when it should have been %s", av.Action, addNode)
	}

	given_hash := av.GetHash()
	av.SetHash([HashSize]byte{0})
	measured_hash := Hash(av)
	av.SetHash(given_hash)
	if av.Hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, av.Hash)
	}

	if av.AdminValid != admin_valid {
		t.Errorf("Expected PrivSol %x, but got %x", admin_valid, av.AdminValid)
	}
	if av.NewNodeID != node {
		t.Errorf("Expected ID %x, but got %x", node, av.NewNodeID)
	}
}
