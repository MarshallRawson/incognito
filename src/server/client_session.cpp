#include "server/client_session.hpp"

ClientSession::ClientSession(int _conn_fd)
  : conn_fd_(_conn_fd)
{}

void
ClientSession::Launch()
{
  std::cout << "client session launched\n";
  // TODO run in own thread
  if (conn_fd_ < 0) {
    throw std::runtime_error("accept failed");
  }
  char client_pubkey[1024] = {};
  int read_ret = read(conn_fd_, client_pubkey, 1024);
  if (read_ret < 0) {
    throw std::runtime_error("read failed");
  }
  e_rsa_.FromPEMStr(client_pubkey);
  std::string server_pubkey = "";
  d_rsa_.PubKeyAsPEMStr(server_pubkey);
  if (send(conn_fd_, server_pubkey.c_str(), strlen(server_pubkey.c_str()), 0) <
      0) {
    throw std::runtime_error("send failed");
  }
  // read encrypted private key
  // decrypt private key

  if (close(conn_fd_)) {
    throw std::runtime_error("Failed to close conn_fd_");
  }
}
