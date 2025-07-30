**Edit** 
I haven't created the best name yet so far I called it benchmark but some names I'm considering is probably "GoSpeed"

# Encrypted file transfer benchmark which aims to find Read/Write encryption speeds. 

Insprised by https://github.com/rfjakob/gocryptfs. 

## What it achevices: 

Measures the performance of file encryption using _crypot/aes_ package between 4 different batch sizes. 

In order to test, you'll need the Go compiler "Go 1.22" which is the lastest compiler as of the creation of this repo. 

Then compile the project using "make build". The following should pop up. 


**_./test.out
Welcome to Encrypted File Trasfer
Size (bytes) | Write (MB/s) | Read (MB/s) | Latency (ms)
1048576       | 1172.39      | 1306.34     | 1.13        
10485760      | 2384.95      | 2638.73     | 7.65        
104857600     | 2213.00      | 3886.96     | 67.48       
1048576000    | 1804.51      | 798.46      | 1876.22     
Test Complete! Don't forget to save your scores!**

If you want to compile your own flags, you can edit the Makefile or any automation tool such as _Just_

## What do I do if my CPU doesn't support AES?

I haven't tested older cpus without AES instruction. So if you want the best possible results make sure to find out whether your cpu support the AES instruction set.

The oldest cpu I have on hand is a 32bit pentium M cpu from the early 2000s. Since the AES-ni instruction set was introduced back in 2010. 

> AES is a symmetric block cipher that encrypts/decrypts data through several rounds. The new 2010 Intel® Core™ processor family (code name Westmere) includes a set of new instructions, Intel® Advanced Encryption Standard (AES) New Instructions (AES-NI).



