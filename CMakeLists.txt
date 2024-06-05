# cmake_minimum_required(VERSION <specify CMake version here>)
project(GameOfResources_Backend)

set(CMAKE_CXX_STANDARD 23)

add_executable(
        test_file.cpp
        src/test.h
        src/math/mymath.cpp
        src/time/time.cpp
)
