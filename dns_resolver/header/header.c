#include "header.h"

uint8_t* encode_domain(uint8_t *cursor, char *domain)
{
  size_t index = 0;
  uint8_t *length_ptr = cursor;
  cursor++;
  uint8_t len = 0;

  while (1)
  {
    if (domain[index] == '.' || domain[index] == '\0')
    {
      *length_ptr = len;
      if (domain[index] == '\0') break;

      length_ptr = cursor;
      cursor++;
      len = 0;
    }
    else 
    {
      *cursor = domain[index];
      cursor++;
      len++;
    }

    index++;
  }

  *cursor = 0x00;
  cursor++;

  return cursor;
}

char *sanitize_domain(char *domain)
{
  char *scheme = strstr(domain, "://");
  if (scheme != NULL)
  {
    domain = scheme + 3;
  }

  if (strncmp(domain, "www.", 4) == 0)
  {
    domain += 4;
  }

  size_t index = 0;
  while (domain[index] != '\0')
  {
    if (domain[index] == '/' || domain[index] == ':')
    {
      domain[index] = '\0';
      break;
    }
    index++;
  }

  return domain;
}
