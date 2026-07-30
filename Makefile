ASSEMBLER = fasm
ASM_SOURCES = wtk.asm

BUILD_DIR = build/

MFV_OFILES = $(BUILD_DIR)mfv.o $(BUILD_DIR)scanner.o $(BUILD_DIR)proj.o
CRAWLER_LIB = $(BUILD_DIR)libcrawler.so
DNS_R_OFILES = $(BUILD_DIR)dnsr.o $(BUILD_DIR)header.o
AP_LIB = $(BUILD_DIR)libap.so

OBJ_FILES = $(BUILD_DIR)wtkit.o $(MFV_OFILES) $(DNS_R_OFILES)

SO_FILES = $(CRAWLER_LIB) $(AP_LIB)

CC = gcc
CFLAGS = -lpthread -lresolv -no-pie -Wl,-rpath,'$$ORIGIN/' -o

GOC = go
GOFLAGS = build -buildmode=c-shared -o 

DB_DIR = ~/.local/share/wtk

build_and_link:
	@mkdir -p $(BUILD_DIR)
	cd mfv && make build_o
	cd crawler && $(GOC) $(GOFLAGS) ../$(CRAWLER_LIB) .
	cd dns_resolver && make build_o
	cd assets_profiler && $(GOC) $(GOFLAGS) ../$(AP_LIB) .
	$(ASSEMBLER) $(ASM_SOURCES) $(BUILD_DIR)wtkit.o
	$(CC) $(OBJ_FILES) $(SO_FILES) $(CFLAGS) $(BUILD_DIR)wtk
	mkdir -p $(DB_DIR)/db

clean:
	rm $(OBJ_FILES) $(BUILD_DIR)wtk
