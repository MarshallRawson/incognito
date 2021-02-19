#include <arpa/inet.h>
#include <bits/stdc++.h>
#include <chrono>
#include <fstream>
#include <iostream>
#include <stdio.h>
#include <string.h>
#include <sys/socket.h>
#include <thread>
#include <unistd.h>

#include "back_end/back_end.hpp"

BackEnd::BackEnd()
{
  InitVerbs();
}

BackEnd::~BackEnd() {}

void
BackEnd::InitVerbs()
{
  verbs_["connect_tcp"] =
    std::bind(&BackEnd::connect_tcp, this, std::placeholders::_1);
}

void
BackEnd::Launch()
{
  while (true) {
    // recive a message from the python front end
    std::string msg = "";
    getline(std::cin, msg);
    std::string verb = msg.substr(0, msg.find(" "));
    if (verbs_.find(verb) != verbs_.end()) {
      // TODO run in own thread
      verbs_[verb](msg);
    } else {
      std::cerr << "NOT RECOGNZED COMMAND: " + msg + "\n";
      exit(1);
    }
  }
}

void
BackEnd::connect_tcp(const std::string& msg)
{
  char* msg_c = new char[strlen(msg.c_str()) + 1];
  strcpy(msg_c, msg.c_str());
  // TODO: add hash to prefix each message
  std::cout << "in_progress\n";
  strtok(msg_c, " ");
  char* ip = strtok(NULL, " ");
  if (ip == NULL) {
    std::cout << "failed \"Malformed Message. Usage: connect <ip> <port>\"\n";
    delete[] msg_c;
    return;
  }
  char* port_str = strtok(NULL, " ");
  if (port_str == NULL) {
    std::cout << "failed \"Malformed Message. Usage: connect <ip> <port>\"\n";
    delete[] msg_c;
    return;
  }
  int port = atoi(port_str);
  tcp_client_.Connect(ip, port);
  std::cout << "success\n";
  delete[] msg_c;
  return;
}
