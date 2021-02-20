#pragma once
#include <stdexcept>
#include <string.h>
#include <string>
#include <unistd.h>
#include <vector>

#include <openssl/pem.h>
#include <openssl/rand.h>
#include <openssl/rsa.h>

class __rsa__
{
public:
  virtual ~__rsa__() = 0;
  void PubKeyAsPEMStr(std::string& _out);
  static const int num_key_bits = 2048;
  static const int padding = RSA_PKCS1_PADDING;
  int Size();

protected:
  __rsa__() {}
  RSA* rsa_ = nullptr;
};
inline __rsa__::~__rsa__() {}

class DecryptionRSA : public __rsa__
{
public:
  DecryptionRSA();
  ~DecryptionRSA();
  std::string Decrypt(const std::vector<unsigned char>& _encryted);

private:
  BIGNUM* bn_ = nullptr;
};

class EncryptionRSA : public __rsa__
{
public:
  void FromPEMStr(const char* _pem);

  ~EncryptionRSA();
  EncryptionRSA(RSA* _rsa);
  EncryptionRSA();
  std::vector<unsigned char> Encrypt(const std::string& _msg);

private:
  void SetRSA(RSA* _rsa);
};
