#include "encryptUtil.h"

int main(int argc, char *argv[]) 
{  
  // copy all arguments in a container of strings
  vector<string> allArgs(argv, argv + argc);
 
	// argc should be 5 for a correct execution
	if (allArgs.size() != 5) {
    // argv[0] is the program name
    cerr << "usage: "<< allArgs[0] << " [-n #] [-k keyfile]\n";
		return -1;
	}
    
  EncryptUtil crypto;
  if(!crypto.read_key_from_file(allArgs[4])) {
  	return -1;
  }
  if(!crypto.set_number_of_threads(allArgs[2])) {
  	return -1;
  }
  crypto.set_stream_buffer_size(10000); // 3 fails
  if(!crypto.input_stream_encrypt()) {
  	return -1;
  }
  
	return 0;
}
