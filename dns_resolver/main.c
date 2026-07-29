#include <sys/socket.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/time.h>
#include "header.h"

int dnsr_main(int argc, char *argv[])
{

  if (argc < 2)
  {
    printf("Usage: %s <url | ip_address>\n", argv[0]);
    exit(EXIT_FAILURE);
  }

  if (strcmp(argv[1], "--help") == 0)
  {
    printf("Usage: %s <url | ip_address>\n", argv[0]);
    exit(EXIT_FAILURE);
  }

  uint8_t buffer[512];
  uint8_t *cursor = buffer;

  PacketHeader *header = (PacketHeader *)cursor;
  header->transaction_id = htons(0x1234);
  header->flags = htons(0x0100);
  header->question_count = htons(1);
  header->answer_count = 0;
  header->authority_count = 0;
  header->additional_count = 0;

  cursor += sizeof(PacketHeader);

  char *domain = sanitize_domain(argv[1]);
  cursor = encode_domain(cursor, domain);

  uint16_t qtype = htons(1);
  memcpy(cursor, &qtype, sizeof(uint16_t));
  cursor += sizeof(uint16_t);

  uint16_t qclass = htons(1);
  memcpy(cursor, &qclass, sizeof(uint16_t));
  cursor += sizeof(uint16_t);

  size_t packet_size = cursor - buffer;

  int socket_fd = socket(AF_INET, SOCK_DGRAM, 0);

  if (socket_fd == -1)
  {
    perror("main: socket");
    exit(EXIT_FAILURE);
  }

  struct sockaddr_in dest_addr;
  memset(&dest_addr, 0, sizeof(struct sockaddr_in));

  dest_addr.sin_family = AF_INET;
  dest_addr.sin_port = htons(53);
  if (inet_pton(AF_INET, "8.8.8.8", &dest_addr.sin_addr) <= 0)
  {
    perror("main: inet_pton");
    close(socket_fd);
    exit(EXIT_FAILURE);
  }

  struct timeval tv;
  tv.tv_sec = 2;
  tv.tv_usec = 0;

  setsockopt(socket_fd, SOL_SOCKET, SO_RCVTIMEO, (const char *)&tv, sizeof(tv));

  ssize_t bytes_sent = sendto(socket_fd, buffer, packet_size, 0, (struct sockaddr *)&dest_addr, sizeof(dest_addr));
  if (bytes_sent == -1)
  {
    perror("main: sendto");
    close(socket_fd);
    exit(EXIT_FAILURE);
  }

  uint8_t response[512];
  struct sockaddr_in server_addr;
  socklen_t addr_len = sizeof(server_addr);

  ssize_t bytes_received = recvfrom(socket_fd, response, sizeof(response), 0, (struct sockaddr *)&server_addr, &addr_len);
  if (bytes_received == -1)
  {
    perror("main: recvfrom");
    close(socket_fd);
    exit(EXIT_FAILURE);
  }

  close(socket_fd);

  uint8_t *response_cursor = response;
  PacketHeader *response_header = (PacketHeader *)response_cursor;
  uint16_t answer_count = ntohs(response_header->answer_count);
  uint16_t response_flags = ntohs(response_header->flags);

  if ((response_flags & 0x000F) != 0)
  {
    printf("DNS Error! Server returned: %d\n", response_flags & 0x000F);
    exit(EXIT_FAILURE);
  }
  if (answer_count == 0)
  {
    printf("No answers found in response.\n");
    exit(EXIT_FAILURE);
  }

  response_cursor += sizeof(PacketHeader);
  while (*response_cursor != 0x00)
  {
    if ((*response_cursor & 0xC0) == 0xC0)
    {
      response_cursor += 2;
      break;
    }
    response_cursor += (*response_cursor + 1);
  }
  if (*response_cursor == 0x00)
  {
    response_cursor++;
  }

  response_cursor += 4;

  for (size_t i = 0; i < answer_count; i++)
  {
    if ((*response_cursor & 0xC0) == 0xC0)
    {
      response_cursor += 2;
    }
    else
    {
      while (*response_cursor != 0x00)
      {
        response_cursor += (*response_cursor + 1);
      }
      response_cursor++;
    }

    uint16_t type = ntohs(*(uint16_t *)response_cursor);
    response_cursor += 2;
    uint16_t class = ntohs(*(uint16_t *)response_cursor);
    response_cursor += 2;
    uint32_t ttl = ntohl(*(uint32_t *)response_cursor);
    response_cursor += 4;
    uint16_t rdlength = ntohs(*(uint16_t *)response_cursor);
    response_cursor += 2;

    if (type == 1 && rdlength == 4)
    {
      char ip_str[INET_ADDRSTRLEN];
      inet_ntop(AF_INET, response_cursor, ip_str, sizeof(ip_str));

      printf("Domain: '%s' resolved to IP: %s (TTL: %u sec) (CLASS: %u) (PACKETS_SENT: %ld) (PACKETS_RECEIVED: %ld)\n", domain, ip_str, ttl, class, bytes_sent, bytes_received);
      exit(EXIT_SUCCESS);
    }

    response_cursor += rdlength;
  }

  printf("Unable to find an IP address\n");
  exit(EXIT_FAILURE);
}
