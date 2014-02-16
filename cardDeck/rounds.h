#ifndef rounds_h
#define rounds_h

#include <stdio.h>
#include <math.h>
#include <stdbool.h>
#include <stdlib.h>

typedef struct card {
	int number;
	struct card_t *next;
}card_t;

typedef struct list {
	card_t *front;
	card_t *back;
}list_t;

unsigned int get_number_of_rounds(unsigned int number_of_cards);

#endif
