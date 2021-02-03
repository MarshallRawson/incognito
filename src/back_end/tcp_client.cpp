#include "back_end/tcp_client.hpp"

TcpClient::TcpClient()
{
}

TcpClient::TcpClient(const std::string& _ip, const int _port)
{
  Init(_ip, _port);
}

void TcpClient::Init(const std::string& _ip, const int _port)
{
  ip_ = _ip;
  port_ = _port;
  serv_addr_.sin_family = AF_INET;
  serv_addr_.sin_port = htons(_port);
  if(inet_pton(AF_INET, ip_.c_str(), &serv_addr_.sin_addr)<=0)
  {
    throw std::runtime_error("Ip / port invalid. ip:" + _ip +
                             " port:" + std::to_string(_port));
  }
  init_ = true;
}

TcpClient::~TcpClient()
{
}

void TcpClient::Connect(const std::string& _ip, const int _port)
{
  Init(_ip, _port);
  Connect();
}

void TcpClient::Connect()
{
  if (init_ == false)
    throw std::runtime_error("Not inited");

  sock_fd_ = socket(AF_INET, SOCK_STREAM, 0);
  if (sock_fd_ < 0)
  {
    throw std::runtime_error("Could not open socket, returned socket fd is:" +
                        std::to_string(sock_fd_) + "\n");
  }
  if (connect(sock_fd_, (struct sockaddr *)&serv_addr_, sizeof(serv_addr_)) < 0)
  {
    throw std::runtime_error("TCP Connection Failed");
  }

  std::string pub_as_str = "";
  d_rsa_.PubKeyAsPEMStr(pub_as_str);

  if (send(sock_fd_, pub_as_str.c_str(), strlen(pub_as_str.c_str()), 0) < 0)
  {
    throw std::runtime_error("Error sending TCP");
  }
  unsigned char buf[1024] = {};
  int read_len = read(sock_fd_, buf, 1024);
  if (read_len < 0)
  {
    throw std::runtime_error("Error reading Server's public key\n");
  }
  e_rsa_.FromPEMStr((char*)buf);
  // Generate a private key
  // encrypt private key
  // send private key
  // use username and a hash of 
  return;
}



