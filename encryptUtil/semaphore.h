#ifndef semaphore_h
#define semaphore_h

#include <mutex>
#include <condition_variable>

class semaphore {
private:
  std::mutex mtx;
  std::condition_variable cv;
  unsigned int count;
    
public:
  semaphore();
  void reset();
  void notify();
  void wait(unsigned int thread_id);
};

#endif
