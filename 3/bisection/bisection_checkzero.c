#include <stdio.h>
#include <math.h>

float bisec(float (*fn)(float), float x_first, float x_last, float delta, float epsilon, float max_iteration) {
    float x_1 = x_first;
    float x_2 = x_first + delta;
    float result;
    int i;

    for (i = 0; i < max_iteration; i++) {
        float v_1 = (*fn)(x_1);
        float v_2 = (*fn)(x_2);
        if (fabs(v_1) < epsilon) return x_1;
        if (fabs(v_2) < epsilon) return x_2;
        if ((result = v_1 * v_2) < 0) {
            if (-result < epsilon)
                return x_2;
            delta = 0.5 * delta;
            bisec(fn, x_1, x_2, delta, epsilon, max_iteration);
        } else {
            x_1 += delta;
            x_2 += delta;
        }
        if (x_1 >= x_last) {
            return x_last + 10;
        }
    }

    return x_last + 5;
}

float func (float x) {
    return x;
}

int main() {
    float x_1, x_2, delta, epsilon, result, start;
    float (*fn)(float) = func;
    int max_iteration = 1000;
    int no_root = 1;

    start = -10;
    x_1 = start;
    x_2 = 10;
    delta = 1;
    epsilon = 0.01;

    do {
        result = bisec(fn, x_1, x_2, delta, epsilon, max_iteration);
        if (result > x_2) {
            printf("There is no more ROOT!!\n");
        } else {
            printf("The %3dth ROOT is %7.3f\n", no_root++, result);
            x_1 = result + delta;
        }
    } while (result < x_2);
}