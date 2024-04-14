#include <stdio.h>

#include "libcomandy.h"

int main() {
    char* msg = request("{}");
    puts(msg);
}
