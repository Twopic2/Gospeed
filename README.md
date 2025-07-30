# Gospeed 
A high-performance encrypted file transfer benchmark written in Go that measures AES encryption/decryption speeds across different file sizes using concurrent operations.

## Overview 
Gospeed is designed to benchmark the performance of AES-256-GCM encryption for file operations. It tests read/write speeds and latency across multiple data sizes, utilizing disk IO. 

** Inspired by gocryptfs **

## Prerequisites 

Go 1.22 or later (latest version recommended)
CPU with AES-NI instruction set support _non-aes or softwware based encryption methods will be implamented in the future._ 
Needs +16Gbs of memory _Femboys loves ram_

## Installation

### Manual Build
```bash
go build -o gospeed .
go run gospeed
```

## Architecture

### Encryption Details
- *Algorithm*: AES-256-GCM (Galois/Counter Mode)
- *Key Size*: 256-bit (32 bytes)
- *Symetrical key exchange*
