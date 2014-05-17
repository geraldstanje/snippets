all: encryptUtil
OBJS = semaphore.o encryptUtil.o main.o
CC = g++
DEBUG = -g
CFLAGS = -Wall -O3 -c $(DEBUG) -std=c++11 -pthread -D_GLIBCXX_USE_NANOSLEEP -D_GLIBCXX_USE_SCHED_YIELD
LFLAGS = -Wall $(DEBUG) -lpthread

encryptUtil : $(OBJS)
	$(CC) $(OBJS) -o encryptUtil $(LFLAGS)

semaphore.o : semaphore.h
	$(CC) $(CFLAGS) semaphore.cpp
	
encryptUtil.o : encryptUtil.h
	$(CC) $(CFLAGS) encryptUtil.cpp

main.o : encryptUtil.h
	$(CC) $(CFLAGS) main.cpp
	    
clean:
	rm -f *.o encryptUtil
