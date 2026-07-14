#include "proj.h"

Project init_project()
{
  Project proj;
  proj.file_count = 0;
  proj.file_capacity = FILE_CAPACITY;
  proj.files = malloc(proj.file_capacity * sizeof(File));
  if (proj.files == NULL)
  {
    perror("init_project: malloc");
    exit(EXIT_FAILURE);
  }

  proj.asset_count = 0;
  proj.asset_capacity = FILE_CAPACITY;
  proj.discovered_assets = malloc(proj.asset_capacity * sizeof(char *));
  if (proj.discovered_assets == NULL)
  {
    perror("init_project: malloc");
    exit(EXIT_FAILURE);
  }

  return proj;
}

void register_physical_asset(Project *proj, const char *path)
{
  if (proj->asset_count >= proj->asset_capacity)
  {
    proj->asset_capacity = proj->asset_capacity == 0 ? FILE_CAPACITY : proj->asset_capacity * 2;

    char **temp_discovered_assets = realloc(proj->discovered_assets, proj->asset_capacity * sizeof(char *));
    if (temp_discovered_assets == NULL)
    {
      perror("register_physical_asset: realloc");
      cleanup_project(proj);
      exit(EXIT_FAILURE);
    }
    proj->discovered_assets = temp_discovered_assets;
  }

  const char *filename = get_filename(path);
  char *lower_name = strdup(filename);
  for (int i = 0; lower_name[i]; i++)
  {
    lower_name[i] = tolower(lower_name[i]);
  }

  proj->discovered_assets[proj->asset_count++] = lower_name;
}

void verify_and_report(Project *proj)
{
  size_t total_missing_files = 0;

  for (size_t i = 0; i < proj->file_count; i++)
  {
    File *file = &proj->files[i];
    bool printed_header = false;

    for (size_t j = 0; file->assets[j].asset_name != NULL; j++)
    {
      AssetFile *ref = &file->assets[j];
      bool found_on_disk = false;

      for (size_t k = 0; k < proj->asset_count; k++)
      {
        if (strcmp(proj->discovered_assets[k], ref->asset_name) == 0)
        {
          found_on_disk = true;
          break;
        }
      }

      if (!found_on_disk)
      {
        if (!printed_header)
        {
          printf(ANSI_BOLD ANSI_COLOR_YELLOW "File Path: %s" ANSI_COLOR_RESET "\n", file->path);
          printed_header = true;
        }
        printf(ANSI_BOLD ANSI_COLOR_RED "\tNot Found: " ANSI_COLOR_RESET "%s\n", ref->path);
        total_missing_files++;
      }
    }
    if (printed_header)
    {
      printf("\n\n");
    }
  }

  printf("\nTotal files missing: %zu\n", total_missing_files);
}

const char *get_filename(const char *path)
{
  const char *last_slash = strrchr(path, '/');
  return (last_slash != NULL) ? last_slash + 1 : path;
}

void cleanup_project(Project *project)
{
  for (size_t i = 0; i < project->file_count; i++)
  {
    free(project->files[i].path);
    free(project->files[i].file_name);
    for (size_t j = 0; project->files[i].assets[j].asset_name != NULL; j++)
    {
      free(project->files[i].assets[j].file_extension);
      free(project->files[i].assets[j].asset_name);
      free(project->files[i].assets[j].path);
    }
    free(project->files[i].assets);
  }
  free(project->files);

  for (size_t i = 0; i < project->asset_count; i++)
  {
    free(project->discovered_assets[i]);
  }
  free(project->discovered_assets);
}
