#include "server/client_com_map.hpp"





int
main (int argc, char** argv)
{
  if (argc < 2)
    exit(1);
  ClientComMap ccm(argv[1]);
  uint64_t alice = 37;
  // since alice is first she can add herself
  ccm.AddClient(alice, alice);
  ccm.LogOn(alice);
  uint64_t bob = 73;
  // alice can now adds bob!
  ccm.AddClient(alice, bob);
  ccm.LogOn(bob);

  std::string hi_bob = "Hi Bob!";
  ccm.PutMsg(alice, hi_bob.c_str(), hi_bob.length()+1);
  ccm.SendMsgTo(alice, bob);

  std::vector<char> msg = ccm.GetMsg(bob);
  std::cout << std::string(&(msg[0])) << "\n";

  ccm.LogOff(alice);
  


  ccm.LogOff(alice);
  ccm.LogOff(bob);
}
