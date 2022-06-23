#include <stdio.h>
#include <stdlib.h>
#include <math.h>

int main() {
    int i, j, k;
    int a[3][3] = {1, 0, 1, 1, 1, 0, 0, 1, 1};
    double c[3] = {0, 90, 80};
    double x[3], factor, sum;

    for (k = 0; k < 3 - 1; k++) {
        for (i = k + 1; i < 3; i++) {
            factor = a[i][k] / a[k][k];
            for (j = k + 1; j < 3; j++)
                a[i][j] -= factor * a[k][j];
            c[i] -= factor * c[k];
        }
    }
    x[3 - 1] = c[3 - 1] / a[3 - 1][3 - 1];
    for (i = 3 - 2; i >= 0; i--) {
        sum = 0;
        for (j = i + 1; j < 3; j++)
            sum += a[i][j] * x[j];
        x[i] = (c[i] - sum) / a[i][i];
    }
    for (i = 0; i < 3; i++)
        printf("x[%1d] = %10.4f, ", i, x[i]);
}