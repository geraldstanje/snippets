#include "rounds.h"

int main(int argc, char *argv[]) {
	// argc should be 2 for a correct execution
	if (argc != 2) {
		// argv[0] is the program name
		printf("usage: %s << #\n", argv[0]);
		return -1;
	}

	printf("number of cards: %d, rounds: %d\n", argv[1], get_number_of_rounds(argv[1]));

	return 0;
}
