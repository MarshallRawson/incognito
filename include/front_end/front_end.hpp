#pragma once
#include <string>

class FrontEnd
{
public:
  FrontEnd(const std::string& _python_file);
  void Launch();

private:
  std::string python_file_;
};
