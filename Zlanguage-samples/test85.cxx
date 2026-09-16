#include <array>
#include <numeric>

int calculate_checksum() {
    const std::array<int, 4> values{1, 2, 3, 4};
    return std::accumulate(values.begin(), values.end(), 0);
}
