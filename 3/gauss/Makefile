CC = gcc
CFLAGS = -g -Wall

TARGETS = gauss_original gauss_switch gauss_switch_error gauss_dynamic

all: $(TARGETS)

$(TARGET): $(TARGET).c
	$(CC) $(CFLAGS) $(TARGET).c -o $(TARGET)

clean:
	rm -f $(TARGETS)

.PHONY : all clean