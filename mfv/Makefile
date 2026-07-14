CC = gcc
CFLAGS = -Wall -Wextra -Iproject -Iscanner
SOURCES = main.c project/proj.c scanner/scanner.c
CLIBS = 

BUILD_DIR = ../build/

build_o:
	$(CC) -o $(BUILD_DIR)mfv.o -c $(CFLAGS) main.c
	$(CC) -o $(BUILD_DIR)scanner.o -c $(CFLAGS) scanner/scanner.c
	$(CC) -o $(BUILD_DIR)proj.o -c $(CFLAGS) project/proj.c
