#include "test82.h"

#include <stddef.h>

int c_fixture_name_length(const char *name) {
    if (name == NULL) {
        return 0;
    }

    const char *cursor = name;
    while (*cursor != '\0') {
        cursor += 1;
    }

    return (int)(cursor - name);
}
