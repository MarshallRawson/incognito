CXX= g++ -std=gnu++17 -g -Wall
C= gcc
PYTHON= -Iinclude -I/usr/include/python3.8 -L/usr/lib/python3.8/config-3.8-x86_64-linux-gnu -L/usr/lib  -lcrypt -lpthread -ldl  -lutil -lm -lm -lpython3.8

OPENSSL= -I/usr/include/openssl -lssl -lcrypto
INCLUDES= -I/usr/include -Iinclude
BACK_END= -Lbuild/back_end -lback_end
FRONT_END= -Lbuild/front_end -lfront_end
TCP_CLIENT= -Lbuild/back_end -ltcp_client
TCP_SERVER= -Lbuild/server -ltcp_server
RSA= -Lbuild/rsa -lrsa
PAGE=-Lbuild -lpage
BLOCK_CHAIN=-Lbuild/block_chain -lblock_chain
CLIENT_SESSION= -Lbuild/server -lclient_session
CONFIGPP= -lconfig++

all: build/client_launch/client_launch \
     build/back_end/launch \
     build/server/launch \
     build/block_chain/test

clean:
	rm -rf build/*/*
	make

build/server/launch: \
  src/server/server.cpp \
  build/server/libtcp_server.a \
  build/server/libclient_session.a
	$(CXX) -o build/server/launch src/server/server.cpp $(INCLUDES) $(TCP_SERVER) $(RSA) $(OPENSSL) $(CLIENT_SESSION)

build/back_end/launch: \
  src/back_end/launch.cpp \
  build/back_end/libback_end.a \
  build/back_end/libtcp_client.a
	$(CXX) -o build/back_end/launch src/back_end/launch.cpp $(INCLUDES) $(BACK_END) $(TCP_CLIENT) $(RSA) $(OPENSSL)

build/client_launch/client_launch: \
  src/client_launch/client_launch.cpp \
  build/front_end/libfront_end.a \
  build/back_end/libback_end.a \
  build/back_end/libtcp_client.a
	$(CXX) -o build/client_launch/client_launch src/client_launch/client_launch.cpp \
    $(INCLUDES) $(FRONT_END) $(PYTHON) $(BACK_END) $(TCP_CLIENT) $(RSA) $(OPENSSL)

build/server/libclient_session.a: \
  src/server/client_session.cpp \
  include/server/client_session.hpp
	$(CXX) -c -o build/server/client_session.o src/server/client_session.cpp $(INCLUDES)
	ar rcs build/server/libclient_session.a build/server/client_session.o

build/block_chain/test: \
  src/block_chain/test.cpp \
  build/block_chain/libblock_chain.a
	$(CXX) -o build/block_chain/test src/block_chain/test.cpp $(INCLUDES) $(BLOCK_CHAIN)

build/block_chain/libblock_chain.a: \
  include/block_chain/block_chain.hpp \
  src/block_chain/block_chain.cpp
	mkdir -p build/block_chain
	$(CXX) -c -o build/block_chain/block_chain.o src/block_chain/block_chain.cpp $(INCLUDES)
	ar rcs build/block_chain/libblock_chain.a build/block_chain/block_chain.o

build/page/libpage.a: \
  src/page/page.cpp \
  include/page/page.hpp
	mkdir -p build/page
	$(CXX) -c -o build/page/page.o src/page/page.cpp $(INCLUDES)
	ar rcs build/page/libpage.a build/page/page.o

build/rsa/librsa.a: \
  src/rsa/rsa.cpp \
  include/rsa/rsa.hpp
	mkdir -p build/rsa
	$(CXX) -c -o build/rsa/rsa.o src/rsa/rsa.cpp $(INCLUDES)
	ar rcs build/rsa/librsa.a build/rsa/rsa.o

build/server/libtcp_server.a: \
  src/server/tcp_server.cpp \
  include/server/tcp_server.hpp \
  build/rsa/librsa.a
	mkdir -p build/server
	$(CXX) -c -o build/server/tcp_server.o src/server/tcp_server.cpp $(INCLUDES)
	ar rcs build/server/libtcp_server.a build/server/tcp_server.o

build/back_end/libtcp_client.a: \
  src/back_end/tcp_client.cpp \
  include/back_end/tcp_client.hpp \
  build/rsa/librsa.a
	mkdir -p build/back_end
	$(CXX) -c -o build/back_end/tcp_client.o src/back_end/tcp_client.cpp $(INCLUDES)
	ar rcs build/back_end/libtcp_client.a build/back_end/tcp_client.o

build/front_end/libfront_end.a: \
  src/front_end/front_end.cpp \
  include/front_end/front_end.hpp
	mkdir -p build/front_end
	$(CXX) -c -o build/front_end/front_end.o src/front_end/front_end.cpp $(PYTHON) $(INCLUDE)
	ar rcs build/front_end/libfront_end.a build/front_end/front_end.o

build/back_end/libback_end.a: \
  src/back_end/back_end.cpp \
  include/back_end/back_end.hpp \
  build/back_end/libtcp_client.a
	mkdir -p build/back_end
	$(CXX) -c -o build/back_end/back_end.o src/back_end/back_end.cpp $(INCLUDES)
	ar rcs build/back_end/libback_end.a build/back_end/back_end.o
