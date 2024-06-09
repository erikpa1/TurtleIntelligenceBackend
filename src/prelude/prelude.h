

#include <vector>
#include <unordered_map>
#include <string>
#include <iostream>


template<typename T>
using HashMap = std::unordered_map<std::string, T>;


//Logging

#define Log  std::cout
#define LogI std::cout << "\033[0m \033[94m [Info] "
#define LogE std::cout << "\033[0m \033[91m [Error] "

#define Endl "\033[0m" << std::endl