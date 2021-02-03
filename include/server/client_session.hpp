#pragma once
#include <sys/socket.h>
#include <netinet/in.h>
#include <stdexcept>
#include <iostream>

#include "rsa.hpp"

class ClientSession
{
public:
  ClientSession(int _conn_fd);
  void Launch();
private:
  EncryptionRSA e_rsa_;
  DecryptionRSA d_rsa_;
  int conn_fd_ = -1;
};
