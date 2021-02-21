#include "server/client_com_map.hpp"

ClientComMap::ClientComMap(const std::string& _working_dir) :
  working_dir_(_working_dir)
{}

ClientComMap::~ClientComMap()
{}

void
ClientComMap::SendMsgTo(const uint64_t& _from, const uint64_t& _to)
{
  if (pub_.find(_from) == pub_.end() || sub_.find(_to) == sub_.end()){
    // do not acknowlege that this publisher does not exist
    // TODO, take the message and trash it
    return;
  }
  // TODO
  // pub_[_from][1] > sub_[_to][0]
}

void
ClientComMap::AddClient(const uint64_t& _pub, const uint64_t& _new_pub)
{
}

bool
ClientComMap::LogOn(const uint64_t& _pub)
{
  if (pub_.find(_pub) == pub_.end()){
    return false;
  }
  else
  {
    // TODO
    // open a pipe and place it
  }
  return false;
}

bool
ClientComMap::LogOff(const uint64_t& _pub)
{
  return false;
}
