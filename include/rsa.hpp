#pragma once
#include <string>
#include <string.h>
#include <stdexcept>
#include <unistd.h>

#include <openssl/rsa.h>
#include <openssl/pem.h>
#include <openssl/rand.h>

class __rsa__
{
public:
  void PubKeyAsPEMStr(std::string& _out);
  static const int num_key_bits = 2048;
  static const int padding = RSA_PKCS1_PADDING;
protected:
  __rsa__() {}
  RSA* rsa_ = nullptr;
};


class DecryptionRSA : public __rsa__
{
public:
  DecryptionRSA();
  ~DecryptionRSA();
  void Decrypt(int _flen, unsigned char* _encrypted, unsigned char* _out);
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
  void Encrypt(int _msg_len, unsigned char* _msg, unsigned char* _out);
  int EncryptMsgSize();
private:
  void SetRSA(RSA* _rsa);
};
