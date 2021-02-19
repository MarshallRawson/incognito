#include "back_end/back_end.hpp"
#include "front_end/front_end.hpp"
#include <iostream>
#include <unistd.h>

int
main(int argc, char** argv)
{
  if (argc < 2) {
    std::cerr << "Usage: client_launch <front_end python file>\n";
    exit(1);
  }

  int front_to_back[2] = {};
  int back_to_front[2] = {};
  pipe(front_to_back);
  pipe(back_to_front);
  if (fork() != 0) {
    // run in child process
    dup2(front_to_back[0], STDIN_FILENO);
    dup2(back_to_front[1], STDOUT_FILENO);

    FrontEnd front_end(argv[1]);
    front_end.Launch();
  } else if (fork() != 0) {
    // run in another child process
    dup2(back_to_front[0], STDIN_FILENO);
    dup2(front_to_back[1], STDOUT_FILENO);

    BackEnd back_end;
    back_end.Launch();
  }
}
