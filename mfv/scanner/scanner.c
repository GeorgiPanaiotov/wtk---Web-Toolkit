#include "scanner.h"

char *targets[] = {".png", ".svg", ".gif", ".jpg", ".jpeg", ".jpe"};

static const char *ext_skip[] = {".dll", ".o", ".so", ".exe", ".min.js", ".log", ".zip", ".gz", ".tar", ".map", ".txt", ".json", ".min.css", ".csproj", ".xml", ".plugin", ".xslt"};

char *get_extension(const char *filename)
{
  const char *dot = strrchr(filename, '.');
  if (!dot)
  {
    return strdup("");
  }
  char *ext = strdup(dot);
  for (int i = 0; ext[i]; i++)
  {
    ext[i] = tolower(ext[i]);
  }
  return ext;
}

bool check_is_asset(const char *path)
{
  const char *dot = strrchr(path, '.');
  if (!dot)
  {
    return false;
  }
  size_t target_count = sizeof(targets) / sizeof(targets[0]);
  for (size_t i = 0; i < target_count; i++)
  {
    if (strcasecmp(dot, targets[i]) == 0)
    {
      return true;
    }
  }
  return false;
}

bool should_skip_dir(const char *directory)
{
  return strcmp(directory, ".") == 0 ||
         strcmp(directory, "..") == 0 ||
         strcmp(directory, "node_modules") == 0 ||
         strcmp(directory, "bin") == 0 ||
         strcmp(directory, "obj") == 0 ||
         strcmp(directory, ".git") == 0;
}

bool is_valid_path(const char *path)
{
  size_t len = strlen(path);
  if (len == 0)
  {
    return false;
  }
  for (size_t i = 0; i < len; i++)
  {
    char c = path[i];
    if (!((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '/' || c == '\\' || c == '.' || c == '_' || c == '-'))
    {
      return false;
    }
  }
  return true;
}

char *load_file(int *fd, off_t file_size)
{
  char *buffer = malloc(file_size + 1);
  if (buffer == NULL)
  {
    return NULL;
  }
  buffer[file_size] = '\0';

  ssize_t read_bytes;
  ssize_t total_read = 0;
  while (total_read < file_size)
  {
    read_bytes = read(*fd, buffer + total_read, READ_CHUNK_SIZE);
    if (read_bytes <= 0)
    {
      break;
    }
    total_read += read_bytes;
  }
  if (read_bytes == -1)
  {
    perror("load_file: read");
    free(buffer);
    close(*fd);
    exit(EXIT_FAILURE);
  }
  return buffer;
}

void parse_source_file_references(Project *proj, const char *source_path, int *fd, off_t file_size)
{
  char *buffer = load_file(fd, file_size);
  if (!buffer)
  {
    return;
  }

  if (proj->file_count >= proj->file_capacity)
  {
    proj->file_capacity *= 2;
    proj->files = realloc(proj->files, proj->file_capacity * sizeof(File));
  }

  File *current_file = &proj->files[proj->file_count];
  current_file->path = strdup(source_path);
  current_file->file_name = strdup(get_filename(source_path));

  size_t asset_alloc_count = 0;
  size_t asset_alloc_capacity = 8;
  current_file->assets = malloc(asset_alloc_capacity * sizeof(AssetFile));

  size_t target_count = sizeof(targets) / sizeof(targets[0]);

  for (size_t t = 0; t < target_count; t++)
  {
    char *ptr = buffer;
    size_t ext_len = strlen(targets[t]);

    while ((ptr = strstr(ptr, targets[t])) != NULL)
    {
      char *start = ptr;
      while (start > buffer && *(start - 1) != '"' && *(start - 1) != '\'' && *(start - 1) != '(')
      {
        start--;
      }

      size_t len = (ptr + ext_len) - start;
      char *extracted = malloc(len + 1);
      strncpy(extracted, start, len);
      extracted[len] = '\0';

      size_t trim_len = strlen(extracted);
      while (trim_len > 0 && (extracted[trim_len - 1] == ')' || extracted[trim_len - 1] == ';' ||
                              extracted[trim_len - 1] == '"' || extracted[trim_len - 1] == '\'' ||
                              extracted[trim_len - 1] == '\r' || extracted[trim_len - 1] == '\n' ||
                              extracted[trim_len - 1] == ' ' || extracted[trim_len - 1] == '\\' ||
                              extracted[trim_len - 1] == ':' || extracted[trim_len - 1] == '}'))
      {
        extracted[trim_len - 1] = '\0';
        trim_len--;
      }

      if (strncmp(extracted, "http", 4) != 0 && strncmp(extracted, "//", 2) != 0 && is_valid_path(extracted) && trim_len > ext_len)
      {
        char char_before_ext = extracted[trim_len - ext_len - 1];
        if (char_before_ext != '/' && char_before_ext != '\\' && char_before_ext != '.')
        {
          bool is_dup = false;
          for (size_t d = 0; d < asset_alloc_count; d++)
          {
            if (strcmp(current_file->assets[d].path, extracted) == 0)
            {
              is_dup = true;
              break;
            }
          }

          if (!is_dup)
          {
            if (asset_alloc_count >= asset_alloc_capacity)
            {
              asset_alloc_capacity *= 2;
              current_file->assets = realloc(current_file->assets, asset_alloc_capacity * sizeof(AssetFile));
            }

            AssetFile *new_asset = &current_file->assets[asset_alloc_count];
            new_asset->path = strdup(extracted);

            const char *fname = get_filename(extracted);
            new_asset->asset_name = strdup(fname);
            for (int i = 0; new_asset->asset_name[i]; i++)
              new_asset->asset_name[i] = tolower(new_asset->asset_name[i]);

            new_asset->file_extension = get_extension(fname);
            asset_alloc_count++;
          }
        }
      }
      free(extracted);
      ptr += ext_len;
    }
  }

  current_file->assets = realloc(current_file->assets, (asset_alloc_count + 1) * sizeof(AssetFile));
  current_file->assets[asset_alloc_count].asset_name = NULL;

  if (asset_alloc_count > 0)
  {
    proj->file_count++;
  }
  else
  {
    free(current_file->path);
    free(current_file->file_name);
    free(current_file->assets);
  }

  free(buffer);
}

void mfv_walk(Project *proj, char *path)
{
  int fd = open(path, O_RDONLY);
  if (fd == -1)
  {
    perror("mfv_walk: open");
    exit(EXIT_FAILURE);
  }

  struct stat st;
  if (fstat(fd, &st) == -1)
  {
    close(fd);
    perror("mfv_walk: fstat");
    exit(EXIT_FAILURE);
  }

  if (S_ISDIR(st.st_mode))
  {
    DIR *dir = fdopendir(fd);
    if (!dir)
    {
      close(fd);
      perror("mfv_walk: fopendir");
      exit(EXIT_FAILURE);
    }

    struct dirent *entry;
    while ((entry = readdir(dir)) != NULL)
    {
      if (should_skip_dir(entry->d_name))
      {
        continue;
      }

      size_t len = strlen(path) + strlen(entry->d_name) + 2;
      char *fullpath = malloc(len);
      if (fullpath == NULL)
      {
        perror("mfv_walk: malloc");
        close(fd);
        exit(EXIT_FAILURE);
      }
      snprintf(fullpath, len, "%s/%s", path, entry->d_name);

      mfv_walk(proj, fullpath);
      free(fullpath);
    }
    closedir(dir);
  }
  else if (S_ISREG(st.st_mode))
  {
    size_t skip_size = sizeof(ext_skip) / sizeof(ext_skip[0]);

    for (size_t i = 0; i < skip_size; i++)
    {
      if (strstr(path, ext_skip[i]) != NULL)
      {
        close(fd);
        return;
      }
    }

    if (check_is_asset(path))
    {
      register_physical_asset(proj, path);
      close(fd);
    }
    else
    {
      parse_source_file_references(proj, path, &fd, st.st_size);
      close(fd);
    }
  }
  else
  {
    close(fd);
  }
}