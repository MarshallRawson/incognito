#include "server/tcp_server.hpp"
#include <iostream>

int
main(int argc, char** argv)
{
  if (argc < 2) {
    std::cerr << "Usage: launch <port>\n";
    exit(EXIT_FAILURE);
  }
  int port = atoi(argv[1]);

  TcpServer server(port);
  server.Launch();
}
