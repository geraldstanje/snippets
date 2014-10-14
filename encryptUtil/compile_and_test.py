#!/usr/bin/python

import os
from struct import *
import sys
import time
import matplotlib.pyplot as plt
import subprocess

# List of test cases
# Test case 0: plain text size is not a multiple of key size
# Test case 1: default test given in the question
# Test case 2: plain text size is a multiple of key size
# Test case 3: plain text size is smaller than key size
# Test case 4: plain text size is longer than key size
# Test case 5: plain text size is very long (depends on resize_factor_test_5) -> exterm test
key                   = ["\x11\x11\x11",
                         "\xf0\xf0",
                         "\x01\x01",
                         "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00",
                         "\x01",
                         "\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"]
plain_text            = ["\x00\x00\x00\x00",
                         "\x01\x02\x03\x04\x11\x12\x13\x14",
                         "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00",
                         "\x01",
                         "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00",
                         "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"]
cypher_text_expected  = ["\x11\x11\x11\x22",
                         "\xf1\xf2\xe2\xe5\xd2\xd1\x94\x93",
                         "\x01\x01\x02\x02\x04\04\x08\x08\x10\x10\x20\x20",
                         "\x01",
                         "\x01\x02\x04\x08\x10\x20\x40\x80\x01\x02\x04\x08\x10\x20\x40\x80\x01\x02\x04\x08\x10\x20\x40\x80\x01\x02\x04\x08\x10\x20\x40\x80\x01\x02\x04\x08\x10\x20\x40\x80\x01\x02\x04\x08\x10\x20\x40\x80\x01\x02\x04\x08\x10\x20\x40\x80\x01\x02\x04\x08\x10\x20\x40\x80",
                         "\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x44\x44\x44\x44\x44\x44\x44\x44\x44\x44\x44\x44\x88\x88\x88\x88\x88\x88\x88\x88\x88\x88\x88\x88"]
threads_num_max = 50
resize_plain_text_test_5 = True # if this is set to True -> resizes testcase 5
resize_factor_test_5 = 60000 # resize factor of test case 5

# compile the encryptUtil program
os.system('clear; make clean; make')

# benchmark table
unit_test_benchmark = []

# run the test suite
for i in range(0,len(key)): 
  unit_test_benchmark.append([]);
  
  for thread_num in range(1,threads_num_max+1):
    # this resizes plain_text for unit test 5
    if resize_plain_text_test_5 and i == 5:
      resize_plain_text_test_5 = False      
      plain_text[i] = plain_text[i] * resize_factor_test_5
      cypher_text_expected[i] = cypher_text_expected[i] * resize_factor_test_5

    # Create the plain_text file
    fo = open("plain_text", "wb")
    for j in range(0,len(plain_text[i])):
      fo.write(plain_text[i][j])
    fo.flush()
    fo.close()

    # Create the key file
    fo = open("key", "wb")
    for j in range(0,len(key[i])):
      fo.write(key[i][j])
    fo.flush()
    fo.close()
    
    # print info to stdout
    sys.stdout.write("Unit Test " + str(i) + " with " + str(thread_num) + " thread(s) started...")
    sys.stdout.flush()
    
    # Execute encryptUtil   
    myinput = open('plain_text')
    myoutput = open('cypher_text', 'w')
    pipe = subprocess.Popen(
        ['./encryptUtil', '-n', str(thread_num), '-k', 'key'],
        stdin=myinput, stdout=myoutput, stderr=subprocess.PIPE)
    stdout, stderr = pipe.communicate()
    benchmark_res = stderr.decode('utf-8')  
    pipe.wait()
    myoutput.flush()
    
    # Print benchmarking info
    print benchmark_res + "ms"
    
    # Save execution time
    unit_test_benchmark[i].append((thread_num, benchmark_res))
    
    # Read the cypher_text file
    fo = open("cypher_text", "rb")
    cypher_text = fo.read()
    fo.flush()
    fo.close()
    
    # Create the cypher_text_expected file
    fo = open("cypher_text_expected", "wb")
    for j in range(0,len(cypher_text_expected[i])):
      fo.write(cypher_text_expected[i][j])
    fo.flush()
    fo.close()
        
    # Compare encryptUtil output with cypher_text_expected
    if cypher_text != cypher_text_expected[i]:
      print "failed"
      
      # print the hexdump of cypher_text_expected, for use of analysis
      print "cypher_text_expected:"
      os.system('hexdump -C cypher_text_expected')
      
      # print the hexdump of cypher_text, for use of analysis
      print "cypher_text:"
      os.system('hexdump -C cypher_text')
      
      # in the error case, we exit the program
      sys.exit()
    else:
      print "passed"

# plot
plot_unit_test = 5
x = [i[0] for i in unit_test_benchmark[plot_unit_test]]
y = [i[1] for i in unit_test_benchmark[plot_unit_test]]
 
plt.plot(x, y)
plt.xlabel('number of threads')
plt.ylabel('execution time [ms]')
plt.title('Unit Test ' + str(plot_unit_test))
plt.show()