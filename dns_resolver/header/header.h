#ifndef HEADER_H
#define HEADER_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <arpa/inet.h>

typedef struct __attribute__((packed)) PacketHeader {
  uint16_t transaction_id;
  uint16_t flags;
  uint16_t question_count;
  uint16_t answer_count;
  uint16_t authority_count;
  uint16_t additional_count;
} PacketHeader;

uint8_t *encode_domain(uint8_t *cursor, char *domain);
char *sanitize_domain(char *domain);


#endif
