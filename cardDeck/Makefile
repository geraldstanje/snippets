all: rounds
OBJS = main.o rounds.o
CC = gcc
DEBUG = -g
CFLAGS = -Wall -O3 -c -std=c11
LFLAGS = -Wall $(DEBUG)

rounds : $(OBJS)
	$(CC) $(OBJS) -o rounds $(LFLAGS)

rounds.o : rounds.h
	$(CC) $(CFLAGS) rounds.c

main.o : rounds.h
	$(CC) $(CFLAGS) main.c
	    
clean:
	rm -f *.o rounds
