#ifndef encryptUtil_h
#define encryptUtil_h

#include <vector>
#include <mutex>
#include <iostream>
#include <string>

#define THREADING
//#define BENCHMARKING
//#define DEBUGINFOS
typedef unsigned char BYTE;

class EncryptUtil {
private:
  std::mutex lock; // protects concurrent write access to stream_buffer or key
  std::vector<BYTE> key; // we store the key in memory
  std::vector<std::vector<BYTE>> thread_key; // each thread stores a copy of the key
  long stream_buffer_index;
  long stream_buffer_size; // is a multiple of the chunk_size
  long chunk_size; // represents the size of the key
  unsigned int number_of_threads;
  unsigned int number_of_rotations;
  
private:
  // checks if the string passed only contains digits
  bool is_digits(const std::string &str);
  // helper functions used for debugging, formats the vector into a hex string
  std::string char_to_hex(unsigned char c);
  std::string get_hex_string(const std::vector<BYTE> &vec);
    
  // rotates the vector by 1 to the left
  void rotate_1bit_left(std::vector<BYTE> &array);
  // rotates the vector by number_of_bits to the left
  void rotate_bits_left(std::vector<BYTE> &array, long number_of_bits);
  
  // this function encrypts the data vector with the key starting at first_processed_chunk
  void xor_encrypt_chunk(std::vector<BYTE> &stream_buffer, 
                         long first_proc_chunk, 
                         long number_of_bytes);
  void xor_encrypt_chunk(std::vector<BYTE> &stream_buffer, 
                         long first_proc_chunk, 
                         long number_of_bytes,
                         long number_of_already_proc_bytes, 
                         long thread_id);
 	
  void copy_key_to_each_thread();
 	
  // this function encrypts the stream buffer, therefore the encryption is equally distributed accross
  // the number of threads specified via the command line argument
  bool stream_buffer_encrypt(std::vector<BYTE> &stream_buffer);
  // this function represents the worker thread
  // the worker thread takes the stream_buffer and encrypts the buffer from the first_processed_chunk
  void thread_function(std::vector<BYTE> &stream_buffer, 
                       long first_processed_chunk, 
                       long number_of_bytes, 
                       long number_of_already_proc_bytes, 
                       long thread_id);

public:
  EncryptUtil();
  EncryptUtil(std::initializer_list<BYTE> l);

  // sets the stream buffer size, which is a multiple of the key size
  void set_stream_buffer_size(const long stream_buff_size);
  // reads the key from the file and assigns it to the key vector
  bool read_key_from_file(const std::string filename);
  // sets the number of threads specified via the command line argument
  bool set_number_of_threads(const std::string num_of_threads_str);
  // processes the enryption on the input stream
  bool input_stream_encrypt();
};

#endif
