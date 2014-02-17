#ifndef rounds_h
#define rounds_h

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
//#define DEBUG_INFOS

typedef struct card {
    unsigned int long number;
    struct card_t *next;
}card_t;

typedef struct list {
    card_t *front;
    card_t *back;
}list_t;

bool get_number_of_rounds(unsigned long number_of_cards, unsigned long *rounds);

#endif
