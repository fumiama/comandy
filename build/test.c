#include <stdio.h>

#include "libcomandy.h"

int main() {
    char* para = "{\"method\":\"GET\","
        "\"url\":\"https://i.pximg.net/img-master/img/2012/04/04/21/24/46/26339586_p0_master1200.jpg\","
        "\"headers\":{"
            "\"Referer\":\"https://www.pixiv.net/\","
            "\"User-Agent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0\""
        "}"
    "}";
    fputs(para, stderr); fputs("\n", stderr);
    puts(request(para));
}
