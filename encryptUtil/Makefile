all: encryptUtil
OBJS = main.o encryptUtil.o
CC = g++
DEBUG = -g
CFLAGS = -Wall -O3 -c $(DEBUG) -std=c++11 -pthread
LFLAGS = -Wall $(DEBUG) -lpthread

encryptUtil : $(OBJS)
	$(CC) $(LFLAGS) $(OBJS) -o encryptUtil

encryptUtil.o : encryptUtil.h
	$(CC) $(CFLAGS) encryptUtil.cpp

main.o : encryptUtil.h
	$(CC) $(CFLAGS) main.cpp
	    
clean:
	rm -f *.o encryptUtil
