#include "rsa/rsa.hpp"
#include <iostream>

int
main(int argc, char** argv)
{
  // create a decyption rsa(public and private key set)
  DecryptionRSA d_rsa;

  // get the public key from the key set as a PEM string
  std::string pem;
  d_rsa.PubKeyAsPEMStr(pem);

  // create an encryption rsa(public rsa key component only) from the PEM string
  EncryptionRSA e_rsa;
  e_rsa.FromPEMStr(pem.c_str());

  // read the message from the command line
  std::string msg;
  std::getline(std::cin, msg, '\n');

  // encrypt the message
  std::vector<unsigned char> e_msg = e_rsa.Encrypt(msg);

  // decrypt the message
  std::string d_msg = d_rsa.Decrypt(e_msg);

  // print the encrypted message (should be the same as the message that was
  // read in)
  std::cout << d_msg << "\n";
}
