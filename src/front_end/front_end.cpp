#include "front_end/front_end.hpp"
#include <Python.h>
#include <stdio.h>

FrontEnd::FrontEnd(const std::string& _python_file)
  : python_file_(_python_file)
{}

void
FrontEnd::Launch()
{
  Py_Initialize();
  FILE* fp = _Py_fopen(python_file_.c_str(), "r");
  PyRun_SimpleFile(fp, python_file_.c_str());
  Py_Finalize();
  return;
}
