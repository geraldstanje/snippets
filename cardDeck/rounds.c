#include "rounds.h"
#include <stdio.h>
#include <stdlib.h>

/*
 * creates a new card node
 */
static card_t *new_card(unsigned long number) {
    card_t *node = (card_t*)malloc(sizeof(card_t));
    
    // error checking
    if(node != NULL) {
        node->number = number;
        node->next = NULL;
    }
    
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
static bool add_new_card(list_t *l, unsigned int number) {
    card_t *c = new_card(number);
    
    // error checking
    if(c == NULL) {
        return false;
    }
    
    if (l->front == NULL) {
        l->front = c;
        l->back = c;
    }
    else {
        l->back->next = c;
        l->back = c;
    }
    
    return true;
}

/*
 * prints the complete list of cards
 */
static void print(list_t *l) {
    card_t *curr = l->front;

    printf("hand: ");
    while (curr != NULL) {
        printf("%lu ", curr->number);
        curr = curr->next;
    }
    printf("\n");
}

/*
 * prints the complete list of cards
 */
static bool init_cards(list_t *hand, unsigned long number_of_cards) {
    unsigned long i = 0;
    bool is_valid = true;
    
    for (i = 0; i < number_of_cards; i++) {
        is_valid = add_new_card(hand, i);
        
	if(!is_valid) {
            return false;
        }
    }
	
    return true;
}

/*
 * checks if the list is sorted by increasing numbers
 */
static bool is_sorted(list_t *hand, int number_of_cards) {
    unsigned long curr_number = 0;
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
    card_t *top_card = NULL;
	
    while (hand->front) { // check if cards are available
        top_card = NULL;

        // Take the top card off the hand and set it on the table
        top_card = remove_front(hand);
        insert_front(table, top_card);
        // no cards on table
        if (!top_card) {
            break;
        }

        // Take the next card off the top and put it on the bottom of the hand in your hand
        top_card = remove_front(hand);
        // no cards in hand
        if (!top_card) {
            break;
        }
        insert_back(hand, top_card);
    }
}

/*
 * picks up all cards from the card from the table and assigns it to hand
 */
static void pickup_hand(list_t **hand, list_t **table) {
    list_t *tmp = *hand;
    *hand = *table;
    *table = tmp;
}

/*
 * allocated memory for hand and table and initializes the hand
 */
static bool init_lists(list_t **hand, list_t **table, unsigned long number_of_cards) {
    // allocate memory for hand
    *hand = create_empty_list();
    if(*hand == NULL) {
        return false; // return false if memory allocation for hand fails
    }
	
    // allocate memory for table
    *table = create_empty_list();
    if(*table == NULL) {
        delete_list(hand); // cleanup already allocated memory from before
        return false; // return false if memory allocation for table fails
    }
	
    // initialize cards for hand
    bool is_valid = init_cards(*hand, number_of_cards);
    if(!is_valid) {
        delete_list(hand); // cleanup already allocated memory from before
        delete_list(table); // cleanup already allocated memory from before
        return false;
    }
    
    return true;
}
    
/*
 * determines how many rounds it will take to put a deck back into the original order
 */
bool get_number_of_rounds(unsigned long number_of_cards, unsigned long *rounds) {
    if (number_of_cards <= 2) {
        *rounds = number_of_cards;
        return true;
    }

    // define hand and table (both are a list of cards)
    list_t *hand = NULL;
    list_t *table = NULL;
    
    // allocate memory for hand and table
    // initialize hand
    if(!init_lists(&hand, &table, number_of_cards)) {
        return false;
    }
    
    bool is_equal = false;
    
    do {
        (*rounds)++;
        next_round(hand, table);
        pickup_hand(&hand, &table);

#ifdef DEBUG_INFOS
        print(hand);
#endif

        is_equal = is_sorted(hand, number_of_cards);
    } while (!is_equal);

    delete_list(&hand);
    delete_list(&table);

    return true;
}
