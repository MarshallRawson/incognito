#include "block_chain/block_chain.hpp"


Block::Block(const std::string& _msg,
             const uint64_t _prev_hash,
             const std::string& _author,
             const Action _action) :
  prev_hash_(_prev_hash), msg_(_msg), author_(_author), action_(_action)
{
  hash_ = Hash<std::string>(msg_ + std::to_string(_prev_hash) +
                            author_ + ActionToString(action_));
}

std::string Block::ActionToString(const Action& a)
{
  switch(a)
  {
    case Action::Post:
      return "Post";
    case Action::AddPublisher:
      return "AddPublisher";
    case Action::RemovePublisher:
      return "RemovePublisher";
    case Action::NewBlockChain:
      return "NewBlockChain";
    case Action::ChangeAuthor:
      return "ChangeAuthor";
  }
  return std::to_string((int)(a));
}

void Block::ToString(std::string& _out) const
{
  _out = author_ + " " + ActionToString(action_) + ":\n" +
         msg_ + "\n";
}


void Block::ToBytes(std::ostream& _out) const
{
  _out.write((char*)&prev_hash_, sizeof(uint64_t));
  _out.write(msg_.c_str(), msg_.length()+1);
  _out.write(author_.c_str(), author_.length()+1);
  _out.write((char*)&action_, sizeof(Action));
}

Block Block::FromBytes(std::istream& _in)
{
  uint64_t prev_hash;
  Action action;

  _in.read((char*)&prev_hash, sizeof(uint64_t));

  std::string msg;
  std::getline(_in, msg, '\0');
  std::string author;
  std::getline(_in, author, '\0');

  _in.read((char*)&action, sizeof(Action));

  return Block(msg, prev_hash, author, action);
}

bool Block::operator==(const Block& rhs) const
{
  return prev_hash_ == rhs.prev_hash_ &&
         hash_ == rhs.hash_ &&
         action_ == rhs.action_ &&
         msg_ == rhs.msg_ &&
         author_ == rhs.author_;
}

bool Block::operator!=(const Block& rhs) const
{
  return !operator==(rhs);
}

BlockChain::BlockChain()
{
}

BlockChain::~BlockChain()
{
}

void BlockChain::ToBytes(std::ostream& _out) const
{
  // TODO
  size_t pubs = publishers_.size();
  _out.write((char*)&pubs, sizeof(size_t));
  for (const auto& i : publishers_)
    _out.write((char*)&i, sizeof(uint64_t));
  size_t blocks = chain_.size();
  _out.write((char*)&blocks, sizeof(size_t));
  for (const auto& i : chain_)
  {
    i.ToBytes(_out);
  }
  size_t authors = authors_.size();
  _out.write((char*)&authors, sizeof(size_t));
  for (const auto& i : authors_)
  {
    _out.write((char*)&i.first, sizeof(uint64_t));
    const char* str = i.second.c_str();
    _out.write((char*)str, i.second.length()+1);
  }
}

BlockChain BlockChain::FromBytes(std::istream& _in)
{
  // TODO
  BlockChain bc;
  size_t pubs = 0;
  _in.read((char*)&pubs, sizeof(size_t));
  for (size_t i = 0; i < pubs; ++i)
  {
    uint64_t pub;
    _in.read((char*)&pub, sizeof(uint64_t));
    bc.publishers_.insert(pub);
  }
  size_t blocks = 0;
  _in.read((char*)&blocks, sizeof(size_t));
  for (size_t i = 0; i < blocks; ++i)
  {
    Block b = Block::FromBytes(_in);
    bc.chain_.push_back(b);
    std::string a = "";
    bc.chain_.back().ToString(a);
  }

  size_t authors = 0;
  _in.read((char*)&authors, sizeof(size_t));
  for (size_t i = 0; i < authors; ++i)
  {
    uint64_t publisher;
    _in.read((char*)&publisher, sizeof(uint64_t));
    std::string auth = "";
    std::getline(_in, auth, '\0');
    bc.authors_.emplace(publisher, auth);
  }
  return bc;
}

BlockChain::AddBlockRet BlockChain::New(const std::string& _init_msg,
                                         const uint64_t _publisher,
                                         const std::string& _author)
{
  return AddBlock(Block(_init_msg, 0, _author, Block::Action::NewBlockChain), _publisher);
}

BlockChain::AddBlockRet BlockChain::AddBlock(const Block& _block, const uint64_t& _publisher)
{
  uint64_t p_hash = Hash<uint64_t>(_publisher);
  // genesis block handling
  if (chain_.size() == 0 && _block.action_ == Block::Action::NewBlockChain)
  {
    publishers_.insert(p_hash);
    authors_.emplace(p_hash, _block.author_);
    chain_.push_back(_block);
    return Success;
  }
  else if (chain_.size() > 0)
  {
    // make sure fits into history
    if (_block.prev_hash_ != LastHash())
      return InvalidPrevHash; // means there has been a conflict and we need to reach a consensus

    // make sure this is comming from a verified publisher for this block chain
    if (publishers_.find(p_hash) == publishers_.end())
      return InvalidPublisher; // this means that this block is comming from an imposter

    // make sure this publisher has used this author name before
    if (_block.author_ != authors_[p_hash])
    {
      return InvalidAuthor;
    }

    // make sure the author names are consistent, unless we want to let them switch,
    //  just make sure it is in the block chain!
    if(_block.action_ == Block::Action::ChangeAuthor)
    {
      authors_[p_hash] = _block.msg_;
      chain_.push_back(_block);
      return Success;
    }

    // if this block is an add publisher block
    if (_block.action_ == Block::Action::AddPublisher)
    {
      std::string new_pub_str = _block.msg_.substr(0, _block.msg_.find(" "));
      uint64_t new_pub;
      try
      {
        new_pub = boost::lexical_cast<int>(new_pub_str);
      }
      catch (boost::bad_lexical_cast& e)
      {
        return InvalidMsg;
      }
      chain_.push_back(_block);

      uint64_t new_pub_hash = Hash<uint64_t>(new_pub);

      std::string new_author = _block.msg_.substr(_block.msg_.find(" ") + 1);
      publishers_.insert(new_pub_hash);
      authors_.emplace(new_pub_hash, new_author);
      return Success;
    }
    else if(_block.action_ == Block::Action::RemovePublisher)
    {
      chain_.push_back(_block);
      publishers_.erase(p_hash);
      authors_.erase(p_hash);
      return Success;
    }

    else if(_block.action_ == Block::Action::Post)
    {
      chain_.push_back(_block);
      return Success;
    }
    else
    {
      return InvalidBlockAction;
    }
  }
  // got wrong initial block type
  else
  {
    return NoHistory;
  }
  return Success;
}


uint64_t BlockChain::LastHash() const
{
  return chain_.back().hash_;
}

void BlockChain::LastMsg(std::string& _out) const
{
 chain_.back().ToString(_out);
}

void BlockChain::LastMsgs(std::vector<std::string>& _out, int n) const
{
  auto it = chain_.rbegin();
  for (int i = 0; i < n; ++i)
  {
    _out.insert(_out.begin(), "");
    it->ToString(*_out.begin());
    if (it == chain_.rend())
      break;
    else
      ++it;
  }
}

void BlockChain::AllMsgs(std::vector<std::string>& _out) const
{
  for (auto it = chain_.rbegin(); it != chain_.rend(); ++it)
  {
    _out.insert(_out.begin(), "");
    it->ToString(*_out.begin());
  }
}

Block BlockChain::LastBlock() const
{
  return chain_.back();
}

void BlockChain::LastBlocks(std::vector<Block>& _out, int n) const
{
  int size = std::max((size_t)n, chain_.size());
  _out.reserve(size);
  auto it = chain_.rbegin();
  for (int i = 0; i < n; ++i)
  {
    _out.insert(_out.begin(), *it);
    if (it == chain_.rend())
      break;
    else
      ++it;
  }
}


void BlockChain::GetChain(std::list<Block>& _out)
{
  _out = chain_;
}

void BlockChain::ResolveConflict()
{
}

bool BlockChain::operator==(const BlockChain& rhs) const
{

  if (publishers_ !=rhs.publishers_)
    return false;

  if (chain_ != rhs.chain_)
    return false;

  if (authors_.size() != rhs.authors_.size())
    return false;


  bool authors_match = true;
  for(const auto& i : authors_)
  {
    if (rhs.authors_.find(i.first) == rhs.authors_.end())
    {
      authors_match = false;
      break;
    }
    else if (authors_.at(i.first) != rhs.authors_.at(i.first))
    {
      authors_match = false;
      break;
    }
  }
  return authors_match;
}

bool BlockChain::operator!=(const BlockChain& rhs) const
{
  return !operator==(rhs);
}
