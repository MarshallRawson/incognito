#include <stdint.h>
#include <unordered_map>

namespace incognito_utils {
template<typename T>
static uint64_t
Hash(const T& t)
{
  // TODO make this a real templated hash
  return std::hash<T>{}(t);
}
}
