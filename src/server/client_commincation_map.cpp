#include "server/client_commuincation_map.hpp"



ClientCommunicationMap::ClientCommunicationMap(const std::string& _working_dir) :
  working_dir_(_working_dir)
{}


void
ClientCommunicationMap::SendMsgTo(const uint64_t& _from, const uint64_t& _to)
{
 if (client_pub_fd_.find(_from) == client_pub_fd_




}







