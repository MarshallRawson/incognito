#include "server/client_com_map.hpp"

using namespace incognito_utils;

ClientComMap::ClientComMap(const std::string& _working_dir)
  : working_dir_(_working_dir)
{}

ClientComMap::~ClientComMap()
{
  for (auto& i : pub_)
  {
    close(i.second.write_);
    close(i.second.read_);
  }
  for (auto& i : sub_)
  {
    close(i.second.write_);
    close(i.second.read_);
  }
}

void
ClientComMap::PutMsg(const uint64_t& _pub, const char* _msg, const int _len)
{
  write(pub_[_pub].write_, (void*)_msg, std::max(_len, msg_length_));
}


std::vector<char>
ClientComMap::GetMsg(const uint64_t& _pub)
{
  std::vector<char> msg(msg_length_, '\0');
  read(sub_[Hash(_pub)].read_, (void*)(&(msg[0])), msg_length_);
  return msg;
}


void
ClientComMap::SendMsgTo(const uint64_t& _from, const uint64_t& _to)
{
  if (pub_.find(_from) == pub_.end() || sub_.find(_to) == sub_.end()) {
    // do not acknowlege that this publisher does not exist
    return;
  }
  // pub_[_from] > sub_[_to]
  char buf[msg_length_] = {};
  read(pub_[_from].read_, (void*)buf, msg_length_);
  write(sub_[_to].write_, (void*)buf, msg_length_);
  memset(buf, 0, msg_length_);
}

void
ClientComMap::AddClient(const uint64_t& _pub, const uint64_t& _new_pub)
{
  //std::cout << _pub << " " << _new_pub << "\n";
  if (sub_.size() != 0) {
    if (pub_.find(_pub) == pub_.end() ||
        pub_.find(_new_pub) != pub_.end() ||
        sub_.find(Hash(_new_pub)) != sub_.end()) {
      std::cout << "skipping " << _new_pub << " " <<
        (pub_.find(_pub) == pub_.end()) << " " <<
        (pub_.find(_new_pub) != pub_.end()) << " " <<
        (sub_.find(Hash(_new_pub)) != sub_.end()) << "\n";
      return;
    }
  }
  FILE* f = fopen(std::string(working_dir_ + std::to_string(_new_pub)).c_str(), "w");
  if (f == nullptr) {
    perror(std::string("could not open __LINE__" +
           working_dir_ + std::to_string(_new_pub)).c_str());
  }
  sub_.emplace(Hash(_new_pub), ComLine{-1, fileno(f)});
  //std::cout << sub_[Hash(_new_pub)].read_ << " " << sub_[Hash(_new_pub)].write_ << "\n";
}

void
ClientComMap::LogOn(const uint64_t& _pub)
{
  if (sub_.find(Hash(_pub)) == pub_.end()) {
    return;
  }
  int sub_fd[2] = {};
  pipe(sub_fd);
  //std::cout << sub_[Hash(_pub)].read_ << " " << sub_[Hash(_pub)].write_ << "\n";
  if (dup2(sub_fd[1], sub_[Hash(_pub)].write_) == -1){
    perror(std::string("Could not redirect file to pipe " +
           std::to_string(__LINE__) + std::to_string(sub_[Hash(_pub)].write_)
            ).c_str());
  }
  sub_[Hash(Hash(_pub))].read_ = sub_fd[0];
  // TODO
  // read the contents of the file corresponding to that publisher into that publisher's
  // recv write
  int pub_fd[2] = {};
  pipe(pub_fd);
  pub_.emplace(_pub, ComLine{pub_fd[0], pub_fd[1]});
}

void
ClientComMap::LogOff(const uint64_t& _pub)
{
  if (pub_.find(_pub) == pub_.end()){
    return;
  }
  pub_.erase(_pub);
  if (sub_.find(Hash(_pub)) == sub_.end()){
    return;
  }
  FILE* f = fopen(std::string(working_dir_ + std::to_string(_pub)).c_str(), "w");
  if (f == nullptr) {
    perror(std::string("could not open __LINE__" +
           working_dir_ + std::to_string(Hash(_pub))).c_str());
  }
  close(sub_[Hash(_pub)].read_);
  if (dup2(sub_[Hash(_pub)].write_, fileno(f)) == -1){
    perror(std::string("Could not redirect file to pipe " + std::to_string(__LINE__)).c_str());
  }
}
