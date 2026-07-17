#include "scanner.h"

int mfv_main(int argc, char *argv[])
{
  if (argc < 2)
  {
    printf("Usage: %s <target_directory>\n", argv[0]);
    return 1;
  }

  Project project = init_project();

  mfv_walk(&project, argv[1]);

  printf("Scan Complete\n");
  verify_and_report(&project);
  cleanup_project(&project);
  return 0;
}