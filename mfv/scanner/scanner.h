#ifndef SCANNER_H
#define SCANNER_H

#include <sys/types.h>

#include "proj.h"

#define READ_CHUNK_SIZE 8192

void mfv_walk(Project *proj, char *path);
char *get_extension(const char *filename);
bool check_is_asset(const char *path);
bool should_skip_dir(const char *directory);
bool is_valid_path(const char *path);
char *load_file(int *fd, off_t file_size);
void parse_source_file_references(Project *proj, const char *source_path, int *fd, off_t file_size);

#endif
