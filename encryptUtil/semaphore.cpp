#include "semaphore.h"

semaphore::semaphore(): count(0) {}
       
void semaphore::reset() {
  std::unique_lock<std::mutex> lck(mtx);
  count = 0;
}
    
void semaphore::notify() {
  std::unique_lock<std::mutex> lck(mtx);
  ++count;
  cv.notify_all();
}
    
void semaphore::wait(unsigned int thread_id) {
  std::unique_lock<std::mutex> lck(mtx);
  cv.wait(lck, [this, &thread_id]() { return count == thread_id; });
}
