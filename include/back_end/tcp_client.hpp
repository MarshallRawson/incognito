#pragma once
#include <arpa/inet.h>
#include <iostream>
#include <stdexcept>
#include <stdio.h>
#include <string.h>
#include <string>
#include <sys/socket.h>
#include <unistd.h>

#include "rsa.hpp"

class TcpClient
{
public:
  TcpClient();
  TcpClient(const std::string& _ip, const int _port);
  void Init(const std::string& _ip, const int _port);
  ~TcpClient();
  void Connect();
  void Connect(const std::string& _ip, const int _port);

private:
  bool init_ = false;

  int sock_fd_ = -1;
  struct sockaddr_in serv_addr_;

  std::string ip_ = "";
  int port_ = -1;

  DecryptionRSA d_rsa_;
  EncryptionRSA e_rsa_;
};
