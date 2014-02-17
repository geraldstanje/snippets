#include "rounds.h"
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <limits.h>
#include <ctype.h>
#include <string.h>

/*
 * checks if the char string is a valid postive number of type unsigned long
 */
bool is_pos_numeric(const char *c_str, unsigned long *number) {
    if (c_str == NULL || 
        *c_str == '\0' || 
        isspace(*c_str) || 
        strstr(c_str,"-")) {
        return false;
    }
    
    errno = 0;
    char *endp;
    *number = strtoul(c_str, &endp, 10);
    
    if (*endp != '\0' ||
        endp == c_str || 
        (*number == ULONG_MAX && errno == ERANGE)) { 
    	return false;
    }
    
    return true;
}

int main(int argc, char *argv[]) {
    // argc should be 2 for a correct execution
    if (argc != 2) {
        printf("usage: %s #number_of_cards\n", argv[0]); // argv[0] is the program name
        return -1;
    }

    unsigned long number_of_cards = 0;
    // converts the second command line argument to a unsigned long
    // checks if the converted number is a positive number
    if(!is_pos_numeric(argv[1], &number_of_cards)) {
        printf("invalid number of rounds\n");
        return -1;
    }
    
    unsigned long rounds = 0;
    if(!get_number_of_rounds(number_of_cards, &rounds)) {
        printf("error: get_number_of_rounds, run out of memory or overflow of variable rounds\n");
        return -1; 
    }
	
    printf("number of cards: %lu, rounds: %lu\n", number_of_cards, rounds);
	
    return 0;
}
