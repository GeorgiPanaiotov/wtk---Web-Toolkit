#ifndef PROJ_H
#define PROJ_H

#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <string.h>
#include <dirent.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>
#include <fcntl.h>
#include <ctype.h>

typedef struct AssetFile
{
  char *file_extension;
  char *asset_name;
  char *path;
} AssetFile;

typedef struct File
{
  char *path;
  char *file_name;
  struct AssetFile *assets;
} File;

typedef struct Project
{
  size_t file_capacity;
  size_t file_count;
  File *files;

  char **discovered_assets;
  size_t asset_count;
  size_t asset_capacity;
} Project;

#define ANSI_COLOR_RED "\x1b[31m"
#define ANSI_COLOR_YELLOW "\x1b[33m"
#define ANSI_COLOR_RESET "\x1b[0m"
#define ANSI_BOLD "\x1b[1m"

#define FILE_CAPACITY 256

Project init_project();
void register_physical_asset(Project *proj, const char *path);
void verify_and_report(Project *proj);
const char *get_filename(const char *path);
void cleanup_project(Project *project);

#endif