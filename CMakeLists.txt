cmake_minimum_required(VERSION 3.26)
project(GameOfResources_Backend)

set(CMAKE_CXX_STANDARD 23)



add_executable(
        test_file.cpp
        src/math/mymath.cpp
        src/time/time.cpp
        src/sim/entity.cpp
        src/sim/entity.h

)
