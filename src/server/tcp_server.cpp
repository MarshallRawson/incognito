#include "server/tcp_server.hpp"
#include "server/client_session.hpp"

TcpServer::TcpServer(int _port)
  : port_(_port)
{
  sock_fd_ = socket(AF_INET, SOCK_STREAM, 0);
  if (sock_fd_ == 0) {
    throw std::runtime_error("socket creation failed");
  }
  opt_ = 1;
  if (setsockopt(sock_fd_,
                 SOL_SOCKET,
                 SO_REUSEADDR | SO_REUSEPORT,
                 &opt_,
                 sizeof(opt_)) != 0) {
    throw std::runtime_error("socketopt failed");
  }
  address_.sin_family = AF_INET;
  address_.sin_addr.s_addr = INADDR_ANY;
  address_.sin_port = htons(port_);
  if (bind(sock_fd_, (struct sockaddr*)&address_, sizeof(address_)) < 0) {
    throw std::runtime_error("bind failed");
  }
}

void
TcpServer::Launch()
{
  std::cout << "tcp server launched\n";
  if (listen(sock_fd_, max_clients) < 0) {
    throw std::runtime_error("listen failed");
  }
  addrlen_ = sizeof(address_);
  while (true) {
    int conn_fd =
      accept(sock_fd_, (struct sockaddr*)&address_, (socklen_t*)&addrlen_);
    // TODO: make sure that the PubKey is PEM string
    ClientSession cs(conn_fd);
    cs.Launch();
  }
}
