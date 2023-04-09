#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <time.h>

int loop = 10000000;

int Sum(int a, int b)
{
    return a + b;
}

int64_t now()
{
    struct timespec t;
    clock_gettime(CLOCK_REALTIME, &t);
    return (int64_t)(t.tv_sec) * 1000000000 + (int64_t)(t.tv_nsec);
}

int main()
{
    int64_t t = now();
    for (int i = 0; i < loop; i++)
    {
        int sum = Sum(i, i + 1);
        sum += 1;
    }
    int64_t t2 = now();
    int64_t used = t2 - t;
    printf("[Sum test] time used: %dns, %dns/op\n", used, used / ((int64_t)(loop)));

    t = now();
    for (int i = 0; i < loop; i++)
    {
        char *p = (char *)malloc(1024);
        p[i % 1024] = i % 256;
        free(p);
    }
    t2 = now();
    used = t2 - t;
    printf("[allocator test]  time used: %dns, %dns/op\n", used, used / ((int64_t)(loop)));
}
