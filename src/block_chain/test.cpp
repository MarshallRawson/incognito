#include "block_chain/block_chain.hpp"
#include <stdint.h>
#include <iostream>
#include <sstream>


int main (int argc, char** argv)
{
  uint64_t alice = 37;
  BlockChain bc;
  BlockChain::AddBlockRet ret;

  {
    // invalid initial post
    Block b("My first Post!", bc.LastHash(), "alice", Block::Action::Post);
    ret = bc.AddBlock(b, alice);
    if (ret != BlockChain::AddBlockRet::NoHistory)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }
  {
    // New Block Chain Test
    ret = bc.New("initial message from Alice!", alice, "alice");
    if (ret != BlockChain::AddBlockRet::Success)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // alice makes a valid post
    Block b("My first real Post!", bc.LastHash(), "alice", Block::Action::Post);
    ret = bc.AddBlock(b, alice);
    if (ret != BlockChain::AddBlockRet::Success)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // fake author Test
    Block b("(fake name)!", bc.LastHash(), "NOT alice", Block::Action::Post);
    ret = bc.AddBlock(b, alice);
    if (ret != BlockChain::AddBlockRet::InvalidAuthor)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }


  uint64_t bob = 15;
  {
    // alice must add bob first!

    Block b("I have not been invited!", bc.LastHash(), "alice", Block::Action::Post);

    ret = bc.AddBlock(b, bob);
    if (ret != BlockChain::AddBlockRet::InvalidPublisher)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  uint64_t old_hash = bc.LastHash(); // used in  a couple tests

  {
    // alice now adds bob (invalid msg)
    Block b("fdsfsefsf " + std::to_string(bob) + " bob", bc.LastHash(), "alice", Block::Action::AddPublisher);
    ret = bc.AddBlock(b, alice);
    if (ret != BlockChain::AddBlockRet::InvalidMsg)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // alice now adds bob (deci)
    Block b(std::to_string(bob) + " bob", bc.LastHash(), "alice", Block::Action::AddPublisher);
    ret = bc.AddBlock(b, alice);
    if (ret != BlockChain::AddBlockRet::Success)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // alice must use latest hash always!
    Block b("Alice from the past!", old_hash, "alice", Block::Action::Post);
    ret = bc.AddBlock(b, alice);
    if (ret != BlockChain::AddBlockRet::InvalidPrevHash)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // bob must use latest hash always!
    Block b("Bob from the past!", old_hash, "bob", Block::Action::Post);
    ret = bc.AddBlock(b, alice);
    if (ret != BlockChain::AddBlockRet::InvalidPrevHash)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // must use a valid Action enum!
    Block b("Alice with an undefined action!", bc.LastHash(), "alice", (Block::Action)(37));
    ret = bc.AddBlock(b, alice);
    if (ret != BlockChain::AddBlockRet::InvalidBlockAction)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // bob makes a valid post
    Block b("My first Post!", bc.LastHash(), "bob", Block::Action::Post);
    ret = bc.AddBlock(b, bob);
    if (ret != BlockChain::AddBlockRet::Success)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // bob would like to change his name
    Block b("xXb0bXx37", bc.LastHash(), "bob", Block::Action::ChangeAuthor);
    ret = bc.AddBlock(b, bob);
    if (ret != BlockChain::AddBlockRet::Success)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // bob makes a valid post
    Block b("Im pickel RICK!", bc.LastHash(), "xXb0bXx37", Block::Action::Post);
    ret = bc.AddBlock(b, bob);
    if (ret != BlockChain::AddBlockRet::Success)
    {
      std::cerr << "Failed. line " << __LINE__ << " Received code: " << ret << "\n";
      exit(1);
    }
  }

  {
    // read back the block chain
    std::vector<std::string> msgs;
    bc.AllMsgs(msgs);
    for (const auto& i : msgs)
    {
      std::cout << i << "\n";
    }
  }

  {
    // block to and from byte test
    Block b = bc.LastBlock();
    std::stringstream ss;
    b.ToBytes(ss);
    Block b1 = Block::FromBytes(ss);

    if (b1 != b)
    {
      std::cerr << "Failed. line " << __LINE__ << "\n";
      std::string a, a1;
      b.ToString(a);
      b1.ToString(a1);
      std::cout << a << a1;
      exit(1);
    }
  }

  {
    // block chain to and from byte test
    BlockChain block_chain = bc;
    std::stringstream st;
    block_chain.ToBytes(st);
    BlockChain block_chain1 = BlockChain::FromBytes(st);
    if (block_chain != block_chain1)
    {
      std::cerr << "Failed. line " << __LINE__ << "\n";
      exit(1);
    }
  }
  exit(0);
}
