#pragma once
#include <functional>
#include <map>
#include <stdio.h>
#include <string.h>
#include <string>
#include <unistd.h>

#include <openssl/md5.h>

#include <libconfig.h++>

#include "back_end/tcp_client.hpp"

class BackEnd
{
public:
  BackEnd();
  ~BackEnd();
  void InitVerbs();
  void Launch();

private:
  // verbs
  void connect_tcp(const std::string& msg);

  std::map<std::string, std::function<void(const std::string&)>> verbs_;

  TcpClient tcp_client_;
};
