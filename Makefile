ASSEMBLER = fasm
ASM_SOURCES = wtk.asm

BUILD_DIR = build/

MFV_OFILES = $(BUILD_DIR)mfv.o $(BUILD_DIR)scanner.o $(BUILD_DIR)proj.o
CRAWLER_LIB = $(BUILD_DIR)libcrawler.a

OBJ_FILES = $(BUILD_DIR)wtkit.o $(MFV_OFILES) $(CRAWLER_LIB)

CC = gcc

GOC = go
GOFLAGS = build -buildmode=c-archive -o ../$(CRAWLER_LIB)

DB_DIR = ~/.local/share/wtk

build_and_link:
	@mkdir -p $(BUILD_DIR)
	cd mfv && make build_o
	cd crawler && $(GOC) $(GOFLAGS) .
	$(ASSEMBLER) $(ASM_SOURCES) $(BUILD_DIR)wtkit.o
	$(CC) $(OBJ_FILES) -lpthread -lresolv -no-pie -o $(BUILD_DIR)wtk
	mkdir -p $(DB_DIR)/db

clean:
	rm $(OBJ_FILES) $(BUILD_DIR)wtk