#pragma once

#include <list>
#include <fstream>
#include <string>
#include <set>
#include <map>
#include <stdint.h>
#include <unordered_map>
#include <vector>
#include <stdexcept>
#include <iostream>
#include <sstream>
#include "boost/lexical_cast.hpp"

class __block_chain__
{
public:
  virtual ~__block_chain__() = 0;
protected:
  template<typename T>
  static uint64_t Hash(const T& t)
  {
    // TODO make this a real templated hash
    return std::hash<T>{}(t);
  }
};
inline __block_chain__::~__block_chain__() {}



struct Block : protected __block_chain__
{
  enum class Action {Post=0, // means msg a text message
                     AddPublisher, // means that msg is a string of digits,
                                   // which is the decimal representation of the new publisher
                                   // followed by a ' ' character then the author name of the publisher
                                   // Example: "1234567 bob"
                     RemovePublisher, // msg is a string of digits,
                                      // which is the decimal version of the publisher to be removed
                     NewBlockChain, // msg is the first message in the block chain
                     ChangeAuthor // msg is the new desired author name
                     };
  Block(const std::string& _msg,
        const uint64_t _prev_hash,
        const std::string& _author,
        const Action _action);
  static std::string ActionToString(const Action& a);
  void ToString(std::string& _out) const;

  void ToBytes(std::ostream& _out) const;
  static Block FromBytes(std::istream& _in);
  bool operator==(const Block& rhs) const;
  bool operator!=(const Block& rhs) const;

  uint64_t prev_hash_;
  std::string msg_;
  uint64_t hash_;
  std::string author_;
  Action action_;
};

class BlockChain : protected __block_chain__
{
public:
  BlockChain();
  ~BlockChain();
  void ToBytes(std::ostream& _out) const;
  static BlockChain FromBytes(std::istream& _in);
  enum AddBlockRet {Success=0,
                    InvalidPrevHash,
                    NoHistory,
                    InvalidPublisher,
                    InvalidMsg,
                    InvalidBlockAction,
                    InvalidAuthor};
  AddBlockRet New(const std::string& _init_msg, const uint64_t _publisher, const std::string& _author);
  AddBlockRet AddBlock(const Block& _block, const uint64_t& _publisher);
  uint64_t LastHash() const;
  void LastMsg(std::string& _out) const;
  void LastMsgs(std::vector<std::string>& _out, int _n) const;
  void AllMsgs(std::vector<std::string>& _out) const;
  Block LastBlock() const;
  void LastBlocks(std::vector<Block>& _out, int _n) const;
  void GetChain(std::list<Block>& _out);
  void ResolveConflict(); //TODO
  void BroadCast(); //TODO
  bool operator==(const BlockChain& rhs) const;
  bool operator!=(const BlockChain& rhs) const;
private:
  std::set<uint64_t> publishers_ = {};
  std::list<Block> chain_ = {};
  std::map<uint64_t, std::string> authors_;
};
