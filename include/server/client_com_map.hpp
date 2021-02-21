#pragma once
#include <iostream>
#include <map>
#include <mutex>
#include <stdexcept>
#include <string.h>
#include <unistd.h>
#include "utils/utils.hpp"
#include <vector>

class ClientComMap
{
public:
  ClientComMap(const std::string& _working_dir);
  ~ClientComMap();
  //static ClientComMap FromBytes(std::istream& _in); // TODO
  //void ToBytes(std::ostream& _out);                 // TODO
  void PutMsg(const uint64_t& _pub, const char* _msg, const int _len);
  std::vector<char> GetMsg(const uint64_t& _pub);
  void SendMsgTo(const uint64_t& _from, const uint64_t& _to);
  void AddClient(const uint64_t& _pub, const uint64_t& _new_pub);
  void LogOn(const uint64_t& _pub);
  void LogOff(const uint64_t& _pub);

  const int msg_length_ = 1024;

private:
  std::string working_dir_;
  struct ComLine
  {
    int read_ = -1;
    int write_ = -1;
  };
  std::map<uint64_t, ComLine> pub_ = {};
  std::map<uint64_t, ComLine> sub_ = {};
};
