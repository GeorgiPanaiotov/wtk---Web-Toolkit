ASSEMBLER = fasm
ASM_SOURCES = wtk.asm

BUILD_DIR = build/

MFV_OFILES = $(BUILD_DIR)mfv.o $(BUILD_DIR)scanner.o $(BUILD_DIR)proj.o

OBJ_FILES = $(BUILD_DIR)wtkit.o $(MFV_OFILES)

CC = gcc

build_and_link:
	cd mfv && make build_o
	$(ASSEMBLER) $(ASM_SOURCES) $(BUILD_DIR)wtkit.o
	$(CC) $(OBJ_FILES) -o $(BUILD_DIR)wtk
	cd crawler && go build
	mv crawler/crawler ./build
	mkdir ~/.local/share/wtk
	mkdir ~/.local/share/wtk/db

clean:
	rm $(OBJ_FILES) $(BUILD_DIR)wtk