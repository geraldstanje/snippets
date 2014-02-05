#!/usr/bin/python

import os
from struct import *
import sys
import time

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
threads_num_max = 10
resize_plain_text_test_5 = True # if this is set to True -> resizes testcase 5
resize_factor_test_5 = 20000 # resize factor of test case 5

# compile the encryptUtil program
os.system('clear; make clean; make')

# run the test suite
for i in range(0,len(key)):
	for thread_num in range(1,threads_num_max+1):
		if resize_plain_text_test_5 == True and i == 5:
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
		cmd = 'cat plain_text | ./encryptUtil -n ' + str(thread_num) + ' -k key > cypher_text'
		os.system(cmd)

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
			
	
	
