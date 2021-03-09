package block_chain

import (
	"crypto/rand"
	"github.com/libp2p/go-libp2p-core/peer"
	"testing"
)

// Testing the block chain
func TestBlockChain(t *testing.T) {
	// Alice constructs her block chain
	bc := New(Self{ID: peer.ID([]byte{3, 7}), Name: "alice"})
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
	bob_chain := New(Self{ID: peer.ID([]byte{7, 3}), Name: "bob"})

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
	err = bc.AddNode(bob_chain.ShareId())
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
	if gen.action != genesis {
		t.Errorf("Expected Action %s, but got %s", genesis, gen.action)
	}
	if gen.name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", gen.name)
	}
	if gen.title != "Alice's page" {
		t.Errorf("Expected Title %s, but got %s", "Alice's page", gen.name)
	}
	if gen.nodeID != peer.ID([]byte{3, 7}) {
		t.Errorf("Expected Ver %x, but got %x", peer.ID([]byte{3, 7}), gen.nodeID)
	}
	hash := gen.hash
	b = b.Next()

	p := b.Value.(*Post)
	if p.action != post {
		t.Errorf("Expected Action %s, but got %s", post, p.action)
	}
	if p.name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", p.name)
	}
	if p.msg != "Hello World!" {
		t.Errorf("Expected  %s, but got %s", "Hello World!", p.msg)
	}
	if p.prevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, p.prevHash)
	}
	hash = p.hash
	b = b.Next()

	pub := b.Value.(*AddPublisher)
	if pub.action != addPublisher {
		t.Errorf("Expected Action %s, but got %s", "AddPub", pub.action)
	}
	if pub.name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", pub.name)
	}
	if pub.newName != "bob" {
		t.Errorf("Expected NewName %s, but got %s", "bob", pub.newName)
	}
	if pub.prevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, pub.prevHash)
	}
	if pub.newPublisherPuzzle != bob_chain.SharePubPuzzle() {
		t.Errorf("Expected  pub puzzle %x, but got %x",
			bob_chain.SharePubPuzzle(), pub.newPublisherPuzzle)
	}
	hash = pub.hash
	b = b.Next()

	post2 := b.Value.(*Post)
	if post2.action != post {
		t.Errorf("Expected Action %s, but got %s", post, post2.action)
	}
	if post2.name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", post2.name)
	}
	if post2.msg != "I added Bob" {
		t.Errorf("Expected  %s, but got %s", "I added Bob", post2.msg)
	}
	if post2.prevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, post2.prevHash)
	}
	hash = post2.hash
	b = b.Next()

	av := b.Value.(*AddNode)
	if av.action != addNode {
		t.Errorf("Expected Action %s, but got %s", addNode, av.action)
	}
	if av.name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", av.name)
	}
	if av.newNodeID != bob_chain.ShareId() {
		t.Errorf("Expected NewVer %x, but got %x", bob_chain.ShareId(), av.newNodeID)
	}
	if av.prevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, av.prevHash)
	}
	hash = av.hash
	b = b.Next()

	nc := b.Value.(*ChangeName)
	if nc.action != changeName {
		t.Errorf("Expected Action %s, but got %s", "NameChange", nc.action)
	}
	if nc.name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", nc.name)
	}
	if nc.newName != "Alice" {
		t.Errorf("Expected New Name %s, but got %s", "Alice", nc.newName)
	}
	if nc.prevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, nc.prevHash)
	}
	hash = nc.hash
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
	rand.Read(pub_valid.solution[:])
	rand.Read(pub_valid.nextPuzzle[:])
	admin_valid := PrivValidation{}
	rand.Read(admin_valid.solution[:])
	rand.Read(admin_valid.nextPuzzle[:])

	gen := MakeGenesis(prev_hash, name, pub_valid, admin_valid, title, peer.ID([]byte{3, 7}))

	if gen.action != genesis {
		t.Errorf("MakeGenesis made block action %s, when it should have been %s",
			gen.action, genesis)
	}

	given_hash := gen.GetHash()
	gen.SetHash([HashSize]byte{0})
	measured_hash := Hash(gen)
	gen.SetHash(given_hash)
	if gen.hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, gen.hash)
	}

	idw := string("Intellectual Dark Web")
	if gen.title != idw {
		t.Errorf("Expected Title %s, but got %s", idw, gen.title)
	}
	node := peer.ID([]byte{3, 7})
	if gen.nodeID != node {
		t.Errorf("Expected Id %x, but got %x", node, gen.nodeID)
	}
}

func TestPostBlock(t *testing.T) {
	name := "SpaceDuck"
	var prevHash [HashSize]byte
	rand.Read(prevHash[:])
	var pub_valid PrivValidation
	rand.Read(pub_valid.solution[:])
	rand.Read(pub_valid.nextPuzzle[:])
	p := MakePost(prevHash, name, pub_valid, "Welcome to the Intellecutial Dark Web!")

	if p.action != post {
		t.Errorf("MakePost made block action %s, when it should have been %s", p.action, post)
	}

	given_hash := p.GetHash()
	p.SetHash([HashSize]byte{0})
	measured_hash := Hash(p)
	p.SetHash(given_hash)
	if p.hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, p.hash)
	}

	idw := string("Welcome to the Intellecutial Dark Web!")
	if p.msg != idw {
		t.Errorf("Expected Post %s, but got %s", idw, p.msg)
	}
}

func TestNameChangeBlock(t *testing.T) {
	name := "SpaceDuck"
	var prevHash [HashSize]byte
	rand.Read(prevHash[:])
	var pub_valid PrivValidation
	rand.Read(pub_valid.solution[:])
	rand.Read(pub_valid.nextPuzzle[:])
	change_name := MakeChangeName(prevHash, name, pub_valid, "$paceDuck")

	if change_name.action != changeName {
		t.Errorf("MakeNameChange made block action %s, when it should have been %s",
			change_name.action, changeName)
	}

	given_hash := change_name.GetHash()
	change_name.SetHash([HashSize]byte{0})
	measured_hash := Hash(change_name)
	change_name.SetHash(given_hash)
	if change_name.hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, change_name.hash)
	}

	name = "$paceDuck"
	if change_name.newName != name {
		t.Errorf("Expected Post %s, but got %s", name, change_name.newName)
	}
}

func TestAddPubBlock(t *testing.T) {
	name := "SpaceDuck"
	var prev_hash [HashSize]byte
	rand.Read(prev_hash[:])
	pub_valid := PrivValidation{}
	rand.Read(pub_valid.solution[:])
	rand.Read(pub_valid.nextPuzzle[:])
	admin_valid := PrivValidation{}
	rand.Read(admin_valid.solution[:])
	rand.Read(admin_valid.nextPuzzle[:])

	var new_pub [PuzzleSize]byte
	rand.Read(new_pub[:])

	ap := MakeAddPublisher(prev_hash, name, pub_valid, admin_valid, new_pub, "SpaceGoose")

	if ap.action != addPublisher {
		t.Errorf("MakePost made block action %s, when it should have been %s", ap.action, addPublisher)
	}

	given_hash := ap.GetHash()
	ap.SetHash([HashSize]byte{0})
	measured_hash := Hash(ap)
	ap.SetHash(given_hash)
	if ap.hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, ap.hash)
	}

	if ap.newPublisherPuzzle != new_pub {
		t.Errorf("Expected NewPub %x, but got %x", new_pub, ap.newPublisherPuzzle)
	}
	me := "SpaceGoose"
	if ap.newName != me {
		t.Errorf("Expected NewName %s, but got %s", me, ap.newName)
	}
	if ap.adminValid.solution != admin_valid.solution {
		t.Errorf("Expected PrivSol %x, but got %x", admin_valid.solution, ap.adminValid.solution)
	}
	if ap.adminValid.nextPuzzle != admin_valid.nextPuzzle {
		t.Errorf("Expected PrivSol %x, but got %x", ap.adminValid.nextPuzzle, admin_valid.nextPuzzle)
	}
}

func TestAddNodeBlock(t *testing.T) {
	name := "SpaceDuck"
	var prev_hash [HashSize]byte
	rand.Read(prev_hash[:])
	pub_valid := PrivValidation{}
	rand.Read(pub_valid.solution[:])
	rand.Read(pub_valid.nextPuzzle[:])
	admin_valid := PrivValidation{}
	rand.Read(admin_valid.solution[:])
	rand.Read(admin_valid.nextPuzzle[:])
	node := peer.ID([]byte{3, 7})

	av := MakeAddNode(prev_hash, name, pub_valid, admin_valid, node)

	if av.action != addNode {
		t.Errorf("MakePost made block action %s, when it should have been %s", av.action, addNode)
	}

	given_hash := av.GetHash()
	av.SetHash([HashSize]byte{0})
	measured_hash := Hash(av)
	av.SetHash(given_hash)
	if av.hash != measured_hash {
		t.Errorf("Expected Hash %x, but got %x", measured_hash, av.hash)
	}

	if av.adminValid != admin_valid {
		t.Errorf("Expected PrivSol %x, but got %x", admin_valid, av.adminValid)
	}
	if av.newNodeID != node {
		t.Errorf("Expected Id %x, but got %x", node, av.newNodeID)
	}
}
