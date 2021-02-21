#pragma once
#include <map>

class ClientComMap
{
public:
  ClientComMap();
  ~ClientComMap();
  FromBytes(istream& _in);
  ToBytes(ostream& _out);
  void SendMsgTo(const uint64_t& _from, const uint64_t& _to);
  void AddClient(const uint64_t& _pub, const uint64_t& _new_pub);
  bool LogOn(const uint64_t& _pub);
  bool LogOff(const uint64_t& _pub);
  


private:
  std::string working_dir_ = "";
  std::map<uint64_t, int[2]> pub_ = {};
  std::map<uint64_t, int[2]> sub_ = {};
};
