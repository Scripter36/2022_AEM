CC = gcc
CFLAGS = -g -Wall

TARGETS = bisection_original bisection_checkzero

all: $(TARGETS)

$(TARGET): $(TARGET).c
	$(CC) $(CFLAGS) $(TARGET).c -o $(TARGET)

clean:
	rm -f $(TARGETS)

.PHONY : all clean