#include "encryptUtil.h"
#include <iterator>
#include <sstream>
#include <bitset>
#include <sstream>
#include <algorithm>
#include <initializer_list>
#include <iomanip>
#include <fstream>
#include <thread>
#include <chrono>
#include <sys/stat.h>
#include <cassert>
#include <chrono>

EncryptUtil::EncryptUtil(): stream_buffer_index(0), 
														stream_buffer_size(0), 
														chunk_size(0), 
														number_of_threads(0) {}
  
EncryptUtil::EncryptUtil(std::initializer_list<BYTE> l): key(l), 
																										     stream_buffer_index(0), 
																										     stream_buffer_size(0), 
																										     number_of_threads(0) {
  chunk_size = key.size();
}

std::string char_to_hex(unsigned char c) {
  std::stringstream stream;
  
	stream << std::setfill('0') << std::setw(2) << std::hex << (int)c;
	
	return stream.str();
}


std::string get_hex_string(const std::vector<BYTE> &vec) {
 std::string out;
	
 for(long i = 0; i < vec.size(); i++) {
   out += char_to_hex(vec[i]);
 }
	
 return out;
}

void EncryptUtil::rotate_1bit_left(std::vector<BYTE> &array) {
	BYTE shifted = 0x00;    
  BYTE overflow = (0x80 & array[0]) >> 7;

  for (long i = (array.size() - 1); i >= 0; i--) {	
  	shifted = (array[i] << 1) | overflow;
    overflow = (0x80 & array[i]) >> 7;
    array[i] = shifted;
  }
}

void EncryptUtil::rotate_bits_left(std::vector<BYTE> &array, 
                                   long number_of_bits) {
	
	number_of_bits = number_of_bits % (array.size() * 8);

  for(long i = 0; i < number_of_bits; i++) { 
  	rotate_1bit_left(array);
  }
}

void EncryptUtil::xor_encrypt_chunk(std::vector<BYTE> &stream_buffer, 
                                    long first_processed_byte, 
                                    long number_of_bytes, 
                                    long number_of_already_proc_bytes, 
                                    long thread_id) {
  long key_index = 0;
	
  for (long byte_index = first_processed_byte; byte_index < first_processed_byte + number_of_bytes; byte_index++) {        	  	
    key_index = key_index % chunk_size;
   
    // this lock protects concurrent write access to stream_buffer or key
   	{
		  std::unique_lock<std::mutex> l(lock);   
		
			stream_buffer[byte_index] = stream_buffer[byte_index] ^ thread_key[thread_id][key_index];
		 
		  key_index++;
			number_of_already_proc_bytes++; 

		  if((number_of_already_proc_bytes % (chunk_size)) == 0) {
		  	rotate_bits_left(thread_key[thread_id], 1);
		  }
    }
    
    // allowing other threads to run
    std::this_thread::yield();
  }
}

void EncryptUtil::xor_encrypt_chunk(std::vector<BYTE> &stream_buffer, 
																	  long first_processed_byte, 
																	  long number_of_bytes) {
  long key_index = 0;
	
  for (long byte_index = first_processed_byte; byte_index < first_processed_byte + number_of_bytes; byte_index++) {        	  	
    key_index = key_index % chunk_size;
			
		stream_buffer[byte_index] = stream_buffer[byte_index] ^ key[key_index];   
    key_index++;
		
    if((key_index % (chunk_size)) == 0) {
    	rotate_bits_left(key, 1);
    }
  }
}

bool EncryptUtil::is_digits(const std::string &str) {
  return all_of(str.begin(), str.end(), ::isdigit);  
}

bool EncryptUtil::set_number_of_threads(const std::string num_of_threads_str) {
  std::stringstream ss;
  ss << num_of_threads_str;
  ss >> number_of_threads;
     
  if(!is_digits(num_of_threads_str) || number_of_threads < 1) {
    std::cerr << "\nError: setting number of threads";
    std::cerr << "\nInfo: set number of threads to 1";
    number_of_threads = 1;  
    thread_key.resize(number_of_threads);
    return false;
  }
  
  thread_key.resize(number_of_threads);  
  return true;
}

bool EncryptUtil::read_key_from_file(const std::string filename) {
  std::ifstream file;
	file.open(filename.c_str(), std::ios::in | std::ios::binary);
	struct stat st;
  stat(filename.c_str(), &st);

	if(!file.is_open()) {
	  std::cerr << "keyfile not found" << '\n';
	  return false;
	}
    
  // read the data
	key.resize(static_cast<size_t>(st.st_size));
  file.read(reinterpret_cast<char*>(&key[0]), static_cast<size_t>(st.st_size));
      		
  if(key.size() > 0) {
    chunk_size = key.size();
  }
      
  // close file
	file.close();
	return true;
}

void EncryptUtil::thread_function(std::vector<BYTE> &stream_buffer, 
                                  long first_proc_chunk, 
                                  long number_of_bytes, 
                                  long number_of_already_proc_bytes, 
                                  long thread_id) {
  xor_encrypt_chunk(stream_buffer, first_proc_chunk, number_of_bytes, number_of_already_proc_bytes, thread_id);
}

void EncryptUtil::set_stream_buffer_size(const long stream_buff_size) {
  stream_buffer_index = 0;           
  stream_buffer_size = stream_buff_size * chunk_size;
  
  for(unsigned int i = 0; i < number_of_threads; i++) {
  	thread_key[i].resize(chunk_size);
  }
}

void EncryptUtil::copy_key_to_each_thread() {
  for(unsigned int i = 0; i < number_of_threads; i++) {
  	copy(key.begin(), key.end(), thread_key[i].begin());
  }
}

bool EncryptUtil::stream_buffer_encrypt(std::vector<BYTE> &stream_buffer) {
  if(chunk_size < 1 || stream_buffer.size() < 1) {
    return false;
  }
  
#ifdef THREADING
  
  // copies the key to each thread
  copy_key_to_each_thread();
  
  std::vector<std::thread> workers;
  long number_of_already_proc_bytes = 0;
  long number_of_bytes = stream_buffer.size();
  long number_of_bytes_per_thread = number_of_bytes / number_of_threads; 

  if(number_of_bytes < chunk_size) {
  	number_of_bytes = stream_buffer.size();
  	number_of_bytes_per_thread = number_of_bytes;
  }
   
  // creates as many threads as the user specified
  for(unsigned int thread_id = 0; thread_id < number_of_threads; thread_id++) { 
    long first_chunk_to_process = number_of_bytes_per_thread * thread_id;
  	
  	if(number_of_bytes > 0) {
  		if(thread_id == number_of_threads - 1) {
  			number_of_bytes_per_thread = number_of_bytes; 	
  		}
  	 	   
  	  if(number_of_bytes_per_thread > 0) {
  	  
#ifdef DEBUGINFOS
  	  	std::cerr << "\nthread_id: " << thread_id << ", first_chunk_to_process: " << first_chunk_to_process << ", number_of_bytes_per_thread: " << number_of_bytes_per_thread;
#endif  	  	
    		// start worker threads
    		workers.push_back(std::thread(&EncryptUtil::thread_function, 
    														      this, 
    														      std::ref(stream_buffer), 
    														      first_chunk_to_process, 
    														      number_of_bytes_per_thread, 
    														      number_of_already_proc_bytes, 
    														      thread_id));

    		number_of_bytes -= number_of_bytes_per_thread;    		
    		number_of_already_proc_bytes += number_of_bytes_per_thread;
    		
    		// rotates the key for the next thread
    		if(thread_id < number_of_threads - 1) {
    			rotate_bits_left(thread_key[thread_id + 1], number_of_already_proc_bytes / chunk_size);
    		}
    	}
    }
  }
        
  // wait for worker threads to complete
  for (std::thread &t: workers) {
    if (t.joinable()) {
      t.join();
    }
  }
  
    
  // rotates the key by stream_buffersize
  rotate_bits_left(key, stream_buffer_size / chunk_size);
  
#else

#ifdef DEBUGINFOS  
	std::cerr << "\nfirst_chunk_to_process: " << 0 << ", stream_buffer.size(): " << stream_buffer.size();
#endif

  xor_encrypt_chunk(stream_buffer, 0, stream_buffer.size());
#endif
  
  stream_buffer_index++;
   
  return true;
}
 
bool EncryptUtil::input_stream_encrypt() {
  long total_read = 0;
  
  if(chunk_size < 1 || stream_buffer_size < 1) {
    std::cerr << "\nError: invalid chunk_size or stream_buffer_size";
    return false;
  }
  
#ifdef BENCHMARKING
  auto start = std::chrono::steady_clock::now();
#endif
  
  do {
    std::vector<BYTE> stream_buffer(stream_buffer_size);
    
    std::cin.read(reinterpret_cast<char*>(&stream_buffer[0]), stream_buffer_size);
     
    total_read = std::cin.gcount();
    stream_buffer.resize(total_read);
    
    if(total_read > 0) {
#ifdef DEBUGINFOS  
			std::cerr << "\ntotal_read: " << total_read;
#endif

    	// perform the encryption on the stream buffer  
    	stream_buffer_encrypt(stream_buffer);
     	
      // print stream_buffer std::vector to stdout
      std::copy(stream_buffer.begin(), stream_buffer.end(), std::ostream_iterator<char>(std::cout));
    }
  }while(total_read > 0);
 
#ifdef BENCHMARKING
  auto end = std::chrono::steady_clock::now();
  auto diff = end - start; 
  std::cerr << "\nruntime stream_encrypt: " << std::chrono::duration <double, std::milli> (diff).count() << "ms";
#endif
  
  return true;
}
