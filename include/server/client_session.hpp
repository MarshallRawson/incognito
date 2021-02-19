#pragma once
#include <iostream>
#include <netinet/in.h>
#include <stdexcept>
#include <sys/socket.h>

#include "rsa/rsa.hpp"

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
