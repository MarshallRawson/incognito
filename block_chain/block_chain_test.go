package block_chain

import (
	"crypto/rand"
	"crypto/sha256"
	"github.com/libp2p/go-libp2p-core/peer"
	"testing"
)

// Testing the block chain

func TestBlockChain(t *testing.T) {
	// Alice constructs her block chain
	bc := New(Self{Id: peer.ID([]byte{3, 7}), Name: "alice"})
	// Alice makes her own block chain
	bc.Genesis("Alice's page")
	// Alice makes a post on her block chain
	bc.Post("Hello World!")

	// Alice meets a person named bob
	bob_chain := New(Self{Id: peer.ID([]byte{7, 3}), Name: "bob"})

	bob_puzzle := bob_chain.SharePubPuzzle()
	// Alice acquires a publisher puzzle from Bob, which bob knows how to solve, but Alice does not

	// Alice gets bob's public puzzle and makes a bock which adds bob as a publisher
	bc.AddPub(bob_puzzle, "bob")

	// Alice makes a post
	bc.Post("I added Bob")

	// Bob Shares with Alice his peerID, and Alice adds him as a verifier
	// so he will recive updates
	bc.AddVer(bob_chain.ShareId())

	// Alice changes her name from "alice" to "Alice"
	bc.NameChange("Alice")

	// TODO
	// check the blocks are valid
	chain := bc.ShareChain()
	b := chain.Front()

	gen := b.Value.(Block)
	g_ex, err := gen.FromGenesis()
	if err != nil {
		panic(err)
	}
	if gen.Action != Genesis {
		t.Errorf("Expected Action %s, but got %s", "Genesis", gen.Action)
	}
	if gen.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", gen.Name)
	}
	if g_ex.Title != "Alice's page" {
		t.Errorf("Expected Title %s, but got %s", "Alice's page", gen.Name)
	}
	if g_ex.Ver != peer.ID([]byte{3, 7}) {
		t.Errorf("Expected Ver %x, but got %x", peer.ID([]byte{3, 7}), g_ex.Ver)
	}
	hash := gen.Hash
	b = b.Next()

	post := b.Value.(Block)
	post_ex, err := post.FromPost()
	if err != nil {
		panic(err)
	}
	if post.Action != Post {
		t.Errorf("Expected Action %s, but got %s", "Post", post.Action)
	}
	if post.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", post.Name)
	}
	if post_ex != "Hello World!" {
		t.Errorf("Expected  %s, but got %s", "Hello World!", post_ex)
	}
	if post.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, post.PrevHash)
	}
	hash = post.Hash
	b = b.Next()

	pub := b.Value.(Block)
	pub_ex, err := pub.FromAddPub()
	if err != nil {
		panic(err)
	}
	if pub.Action != AddPub {
		t.Errorf("Expected Action %s, but got %s", "AddPub", pub.Action)
	}
	if post.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", pub.Name)
	}
	if pub_ex.NewName != "bob" {
		t.Errorf("Expected NewName %s, but got %s", "bob", pub_ex.NewName)
	}
	if pub.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, pub.PrevHash)
	}
	if pub_ex.NewPub != bob_chain.SharePubPuzzle() {
		t.Errorf("Expected  pub puzzle %x, but got %x", bob_chain.SharePubPuzzle(), pub_ex.NewPub)
	}
	hash = pub.Hash
	b = b.Next()

	post2 := b.Value.(Block)
	post2_ex, err := post2.FromPost()
	if err != nil {
		panic(err)
	}
	if post2.Action != Post {
		t.Errorf("Expected Action %s, but got %s", "Post", post2.Action)
	}
	if post2.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", post2.Name)
	}
	if post2_ex != "I added Bob" {
		t.Errorf("Expected  %s, but got %s", "I added Bob", post2_ex)
	}
	if post2.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, post2.PrevHash)
	}
	hash = post2.Hash
	b = b.Next()

	av := b.Value.(Block)
	av_ex, err := av.FromAddVer()
	if err != nil {
		panic(err)
	}
	if av.Action != AddVer {
		t.Errorf("Expected Action %s, but got %s", "AddVer", av.Action)
	}
	if av.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", av.Name)
	}
	if av_ex.NewVer != bob_chain.ShareId() {
		t.Errorf("Expected NewVer %x, but got %x", bob_chain.ShareId(), av_ex.NewVer)
	}
	if av.PrevHash != hash {
		t.Errorf("Expected  PrevHash %x, but got %x", hash, av.PrevHash)
	}
	hash = av.Hash
	b = b.Next()

	nc := b.Value.(Block)
	nc_ex, err := nc.FromNameChange()
	if err != nil {
		panic(err)
	}
	if nc.Action != NameChange {
		t.Errorf("Expected Action %s, but got %s", "NameChange", nc.Action)
	}
	if nc.Name != "alice" {
		t.Errorf("Expected Name %s, but got %s", "alice", nc.Name)
	}
	if nc_ex != "Alice" {
		t.Errorf("Expected New Name %s, but got %s", "Alice", nc_ex)
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

func testBlockCoreEqual(a Block, b Block) bool {
	if a.PrevHash != b.PrevHash {
		return false
	}
	if a.PubSol != b.PubSol {
		return false
	}
	if a.NextPub != b.NextPub {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	return true
}

func TestGenesisBlock(t *testing.T) {
	gen := Block{}
	gen.Name = "Spaceduck"
	rand.Read(gen.PrevHash[:])
	rand.Read(gen.PubSol[:])
	rand.Read(gen.NextPub[:])

	var gen_cpy = gen

	var priv [PuzzleSize]byte
	rand.Read(priv[:])
	gen.MakeGenesis("Intellectual Dark Web", priv, peer.ID([]byte{3, 7}))

	if gen.Action != Genesis {
		t.Errorf("MakeGenesis made block action %s, when it should have been %s", gen.Action, Genesis)
	}

	if testBlockCoreEqual(gen, gen_cpy) != true {
		t.Errorf("MakeGenesis Modified the the block core")
	}

	hash := sha256.Sum256(gen.serialize())
	if gen.Hash != hash {
		t.Errorf("Expected Hash %x, but got %x", hash, gen.Hash)
	}

	gen_ex, err := gen.FromGenesis()
	if err != nil {
		panic(err)
	}
	idw := string("Intellectual Dark Web")
	if gen_ex.Title != idw {
		t.Errorf("Expected Title %s, but got %s", idw, gen_ex.Title)
	}
	ver := peer.ID([]byte{3, 7})
	if gen_ex.Ver != ver {
		t.Errorf("Expected Id %x, but got %x", ver, gen_ex.Ver)
	}
}

func TestPostBlock(t *testing.T) {
	post := Block{}
	rand.Read(post.PrevHash[:])
	rand.Read(post.PubSol[:])
	rand.Read(post.NextPub[:])

	var post_cpy = post

	post.MakePost("Welcome to the Intellecutial Dark Web!")

	if post.Action != Post {
		t.Errorf("MakePost made block action %s, when it should have been %s", post.Action, Post)
	}
	if testBlockCoreEqual(post, post_cpy) != true {
		t.Errorf("MakePost Modified the the block core")
	}

	hash := sha256.Sum256(post.serialize())
	if post.Hash != hash {
		t.Errorf("Expected Hash %x, but got %x", hash, post.Hash)
	}

	post_str, err := post.FromPost()
	if err != nil {
		panic(err)
	}

	idw := string("Welcome to the Intellecutial Dark Web!")
	if post_str != idw {
		t.Errorf("Expected Post %s, but got %s", idw, post_str)
	}
}

func TestNameChangeBlock(t *testing.T) {
	post := Block{}
	rand.Read(post.PrevHash[:])
	rand.Read(post.PubSol[:])
	rand.Read(post.NextPub[:])

	var post_cpy = post

	post.MakeNameChange("$paceDuck")

	if post.Action != NameChange {
		t.Errorf("MakeNameChange made block action %s, when it should have been %s", post.Action,
			NameChange)
	}
	if testBlockCoreEqual(post, post_cpy) != true {
		t.Errorf("MakeNameChange Modified the the block core")
	}

	hash := sha256.Sum256(post.serialize())
	if post.Hash != hash {
		t.Errorf("Expected Hash %x, but got %x", hash, post.Hash)
	}

	new_name, err := post.FromNameChange()
	if err != nil {
		panic(err)
	}

	name := string("$paceDuck")
	if new_name != name {
		t.Errorf("Expected Post %s, but got %s", name, new_name)
	}
}

func TestAddPubBlock(t *testing.T) {
	ap := Block{}
	rand.Read(ap.PrevHash[:])
	rand.Read(ap.PubSol[:])
	rand.Read(ap.NextPub[:])

	var ap_cpy = ap

	var new_pub [PuzzleSize]byte
	rand.Read(new_pub[:])
	var priv_sol [SolutionSize]byte
	rand.Read(priv_sol[:])
	var next_priv [PuzzleSize]byte
	rand.Read(next_priv[:])
	ap.MakeAddPub(AddPubExtras{new_pub, "SpaceGoose", priv_sol, next_priv})

	if ap.Action != AddPub {
		t.Errorf("MakePost made block action %s, when it should have been %s", ap.Action, AddPub)
	}
	if testBlockCoreEqual(ap, ap_cpy) != true {
		t.Errorf("MakeAddPub Modified the the block core")
	}

	hash := sha256.Sum256(ap.serialize())
	if ap.Hash != hash {
		t.Errorf("Expected Hash %x, but got %x", hash, ap.Hash)
	}

	ap_ex, err := ap.FromAddPub()
	if err != nil {
		panic(err)
	}
	if ap_ex.NewPub != new_pub {
		t.Errorf("Expected NewPub %x, but got %x", new_pub, ap_ex.NewPub)
	}
	me := string("SpaceGoose")
	if ap_ex.NewName != me {
		t.Errorf("Expected NewName %s, but got %s", me, ap_ex.NewName)
	}
	if ap_ex.PrivSol != priv_sol {
		t.Errorf("Expected PrivSol %x, but got %x", priv_sol, ap_ex.PrivSol)
	}
	if ap_ex.NextPriv != next_priv {
		t.Errorf("Expected PrivSol %x, but got %x", next_priv, ap_ex.NextPriv)
	}
}

func TestAddVerBlock(t *testing.T) {
	av := Block{}
	rand.Read(av.PrevHash[:])
	rand.Read(av.PubSol[:])
	rand.Read(av.NextPub[:])

	var av_cpy = av

	var priv_sol [SolutionSize]byte
	rand.Read(priv_sol[:])
	var next_priv [PuzzleSize]byte
	rand.Read(next_priv[:])
	ver := peer.ID([]byte{3, 7})
	av.MakeAddVer(AddVerExtras{priv_sol, next_priv, ver})

	if av.Action != AddVer {
		t.Errorf("MakePost made block action %s, when it should have been %s", av.Action, AddVer)
	}
	if testBlockCoreEqual(av, av_cpy) != true {
		t.Errorf("MakeAddPub Modified the the block core")
	}

	hash := sha256.Sum256(av.serialize())
	if av.Hash != hash {
		t.Errorf("Expected Hash %x, but got %x", hash, av.Hash)
	}

	av_ex, err := av.FromAddVer()
	if err != nil {
		panic(err)
	}
	if av_ex.PrivSol != priv_sol {
		t.Errorf("Expected PrivSol %x, but got %x", priv_sol, av_ex.PrivSol)
	}
	if av_ex.NextPriv != next_priv {
		t.Errorf("Expected PrivSol %x, but got %x", next_priv, av_ex.NextPriv)
	}
	if av_ex.NewVer != ver {
		t.Errorf("Expected Id %x, but got %x", ver, av_ex.NewVer)
	}
}
