#ifndef encryptUtil_h
#define encryptUtil_h

#include <iostream>
#include <iterator>
#include <sstream>
#include <bitset>
#include <string>
#include <sstream>
#include <vector>
#include <algorithm>
#include <initializer_list>
#include <iomanip>
#include <fstream>
#include <thread>
#include <mutex>
#include <chrono>
#include <sys/stat.h>
#include <cassert>
using namespace std;

#define THREADING
//#define BENCHMARKING
//#define DEBUGINFOS
typedef unsigned char BYTE;

class EncryptUtil {
private:
	mutex m;
	vector<BYTE> key; // we store the key in memory
	vector<vector<BYTE>> thread_key;
	long stream_buffer_index;
	long stream_buffer_size; // is a multiple of the junk_size
	long chunk_size; // represents the size of the key
	unsigned int number_of_threads;

private:
	// checks if the string passed only contains digits
	bool is_digits(const string &str);

	// calculates the number of junks
	long set_num_of_junks(long stream_buffer_size);

	// rotates the key vector by number_of_bits to the left
	void rotate_1bit_left(vector<BYTE> &array);
	void rotate_bits_left(vector<BYTE> &array, long number_of_bits);

	// this function encrypts the data vector with the key starting at first_processed_junk
	void xor_encrypt_junk(vector<BYTE> &stream_buffer, long first_proc_chunk, long number_of_bytes);
	void xor_encrypt_junk(vector<BYTE> &stream_buffer, long first_proc_chunk, long number_of_bytes, long number_of_already_proc_bytes, long thread_id);

	void copy_key_to_each_thread();

	// this function encrypts the stream buffer, therefore the encryption is equally distributed accross
	// the number of threads specified via the command line argument
	bool stream_buffer_encrypt(vector<BYTE> &stream_buffer);
	// this function represents the worker thread
	// the worker thread takes the stream_buffer and encrypts the buffer from the first_processed_junk
	void thread_function(vector<BYTE> &stream_buffer, long first_processed_chunk, long number_of_bytes, long number_of_already_proc_bytes, long thread_id);

public:
	EncryptUtil();
	EncryptUtil(initializer_list<BYTE> l);

	// sets the stream buffer size, which is a multiple of the key size
	void set_stream_buffer_size(const long stream_buff_size);
	// reads the key from the file and assigns it to the key vector
	bool read_key_from_file(const string filename);
	// sets the number of threads specified via the command line argument
	bool set_number_of_threads(string num_of_threads_str);
	// processes the enryption on the input stream
	bool input_stream_encrypt();
};

#endif
