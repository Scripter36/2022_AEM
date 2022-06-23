#include <stdio.h>
#include <stdlib.h>
#include <math.h>

#define MAX_N 10

int main() {
    int i, j, k, N;
    // N을 입력받는다.
    printf("구하고자 하는 미지수의 개수를 입력하세요: ");
    scanf("%d", &N);
    // N의 크기 검증
    if (N > MAX_N) {
        fprintf(stderr, "N의 값이 너무 큽니다.\n");
        exit(EXIT_FAILURE);
    }
    if (N < 1) {
        fprintf(stderr, "N의 값이 너무 작습니다.\n");
        exit(EXIT_FAILURE);
    }
    // 할당
    double *a = malloc(sizeof(double) * N * N);
    double *c = malloc(sizeof(double) * N);
    double *x = malloc(sizeof(double) * N);
    double factor, sum;
    // 입력을 받아 준다.
    printf("왼쪽 행렬을 입력하세요: \n");
    for (i = 0; i < N; i++) {
        for (j = 0; j < N; j++) {
            scanf("%lf", &a[i * N + j]);
        }
    }
    printf("오른쪽 행렬을 입력하세요: \n");
    for (i = 0; i < N; i++) {
        scanf("%lf", &c[i]);
    }

    for (k = 0; k < N - 1; k++) {
        if (a[k * N + k] == 0) { // 대각성분이 0인 경우,
            for (j = k + 1; j < N; j++) {
                if (a[j * N + k] != 0) { // 아래에서 대각성분이 0이 아닌 것을 발견한 경우,
                    // k행과 j행을 바꾼다.
                    double temp;
                    for (i = 0; i < N; i++) {
                        temp = a[j * N + i];
                        a[j * N + i] = a[k * N + i];
                        a[k * N + i] = temp;
                    }
                    break;
                }
                if (j == N - 1) {
                    // 전부 확인했음에도 0이 아닌 경우에만 진입할 수 있다.
                    fprintf(stderr, "해가 존재하지 않습니다.\n");
                    exit(EXIT_FAILURE);
                }
            }
        }
        for (i = k + 1; i < N; i++) {
            factor = a[i * N + k] / a[k * N + k];
            for (j = k; j < N; j++)
                a[i * N + j] -= factor * a[k * N + j];
            c[i] -= factor * c[k];
        }
    }
    x[N - 1] = c[N - 1] / a[(N - 1) * N + N - 1];
    for (i = N - 2; i >= 0; i--) {
        sum = 0;
        for (j = i + 1; j < N; j++)
            sum += a[i * N + j] * x[j];
        x[i] = (c[i] - sum) / a[i * N + i];
    }
    for (i = 0; i < N; i++)
        printf("x[%1d] = %10.4f, ", i, x[i]);
    printf("\n");
}