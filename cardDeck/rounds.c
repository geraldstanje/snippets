#include <stdio.h>
#include <math.h>
#include <stdbool.h>
#include <stdlib.h>

/*
 * creates a new card on the heap
 */
static card_t *new_card(int number) {
	card_t *node = (card_t*)malloc(sizeof(card_t));
	node->number = number;
	node->next = NULL;
	return node;
}

/*
 * allocates an empty list  
 */
static list_t *create_empty_list() {
	list_t *l = (list_t*)malloc(sizeof(list_t));
	l->front = NULL;
	l->back = NULL;

	return l;
}

/*
 * deallocates the entire list
 */
static void delete_list(list_t **l) {
	card_t *curr = (*l)->front;

	while (curr) {
		card_t *to_del = curr;
		curr = curr->next;
		free(to_del);
	}

	free(*l);
	*l = NULL;
}

/*
 * removes the first card from the list and returns it
 */
static card_t *remove_front(list_t *l) {
	card_t *front = l->front;

	if (l->front) {
		l->front = l->front->next;
	}

	return front;
}

/*
 * inserts a new card at the beginning of the list
 */
static void insert_front(list_t *l, card_t *c) {
	c->next = NULL;

	if (l->front == NULL) {
		l->front = c;
		l->back = c;
	}
	else {
		c->next = l->front;
		l->front = c;
	}
}

/*
 * inserts a new card at the end of the list
 */
static void insert_back(list_t *l, card_t *c) {
	c->next = NULL;

	if (l->front == NULL) {
		l->front = c;
		l->back = c;
	}
	else {
		l->back->next = c;
		l->back = c;
	}
}

/*
 * inserts a new card at the end of the list
 */
static void add_new_card(list_t *l, unsigned int number) {
	card_t *c = new_card(number);

	if (l->front == NULL) {
		l->front = c;
		l->back = c;
	}
	else {
		l->back->next = c;
		l->back = c;
	}
}

/*
 * prints the complete list of cards
 */
static void print(list_t *l) {
	card_t *curr = l->front;

	while (curr != NULL) {
		printf("%d ", curr->number);
		curr = curr->next;
	}
	printf("\n");
}

/*
 * prints the complete list of cards
 */
static void init_cards(list_t *hand, unsigned int number_of_cards) {
	unsigned int i = 0;

	for (i = 0; i < number_of_cards; i++) {
		add_new_card(hand, i);
	}
}

/*
 * checks if the list is sorted by increasing numbers
 */
static bool is_sorted(list_t *hand, int number_of_cards) {
	unsigned int curr_number = 0;
	card_t *curr = hand->front;

	for (curr_number = 0; curr_number < number_of_cards; curr_number++) {
		if (curr->number != curr_number) {
			return false;
		}

		curr = curr->next;
	}

	return true;
}

/*
 * runs one round of the card shuffling algorithm
 */
static void next_round(list_t *hand, list_t *table) {
	while (hand->front) {
		card_t *top_card = NULL;

		// Take the top card off the hand and set it on the table
		top_card = remove_front(hand);
		insert_front(table, top_card);

		if (!top_card) {
			break;
		}

		// Take the next card off the top and put it on the bottom of the hand in your hand
		top_card = remove_front(hand);
		if (!top_card) {
			break;
		}
		insert_back(hand, top_card);
	}
}

/*
 * picks up the new hand from the table
 */
static void pickup_hand(list_t **hand, list_t **table) {
	list_t *tmp = *hand;
	*hand = *table;
	*table = tmp;
}

/*
 * determines how many rounds it will take to put a deck back into the original order
 */
unsigned int get_number_of_rounds(unsigned int number_of_cards) {
	if (number_of_cards <= 2) {
		return number_of_cards;
	}

	list_t *hand = create_empty_list();
	list_t *table = create_empty_list();
	bool is_equal = false;
	unsigned int rounds = 0;

	init_cards(hand, number_of_cards);

	do {
		rounds++;
		next_round(hand, table);
		pickup_hand(&hand, &table);
		//print(hand);
		is_equal = is_sorted(hand, number_of_cards);
	} while (!is_equal);

	delete_list(&hand);
	delete_list(&table);

	return rounds;
}