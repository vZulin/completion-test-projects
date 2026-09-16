cmake_minimum_required(VERSION 3.20)
project(cpp_smoke LANGUAGES C CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_executable(cpp_smoke test81.cpp test82.h test83.c test84.cc test85.cxx)
