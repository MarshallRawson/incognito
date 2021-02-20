#pragma once
#include <map>
#include 

class ClientCommuincationMap
{
public:
  ClientCommuincationMap(const std::string& _working_dir);
  ~ClientCommuincationMap();
  void SendMsgTo(const uint64_t& _from, const uint64_t& _to);
  bool LogOn(const uint64_t& pub);
  bool LogOff(const uint64_t& pub);

private:
  std::string working_dir_;
  std::map<uint64_t, FILE*> pub_;
  std::map<uint64_t, FILE*> sub_;
};
