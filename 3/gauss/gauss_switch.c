#include <stdio.h>
#include <stdlib.h>
#include <math.h>

int main() {
    int i, j, k;
    int a[3][3] = {0, 1, 1, 1, 0, 1, 1, 1, 0};
    double c[3] = {2, 2, 2};
    double x[3], factor, sum;

    for (k = 0; k < 3 - 1; k++) {
        if (a[k][k] == 0) { // 대각성분이 0인 경우,
            for (j = k + 1; j < 3; j++) {
                if (a[j][k] != 0) { // 아래에서 대각성분이 0이 아닌 것을 발견한 경우,
                    // k행과 j행을 바꾼다.
                    int temp;
                    for (i = 0; i < 3; i++) {
                        temp = a[j][i];
                        a[j][i] = a[k][i];
                        a[k][i] = temp;
                    }
                    break;
                }
            }
        }
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