#include "server/tcp_server.hpp"

int
main(int argc, char** argv)
{
  if (argc < 2) {
    perror("Usage: launch <port>");
    exit(EXIT_FAILURE);
  }
  int port = atoi(argv[1]);

  TcpServer server(port);
  server.Launch();
}
