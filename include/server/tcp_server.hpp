#pragma once
#include <netinet/in.h>
#include <stdexcept>
#include <sys/socket.h>

#include "rsa.hpp"

class TcpServer
{
public:
  TcpServer(int _port);

  void Launch();

private:
  int port_ = -1;
  int sock_fd_ = -1;
  int opt_ = 1;
  struct sockaddr_in address_;
  int addrlen_ = 0;

  static const int max_clients = 30;
};
