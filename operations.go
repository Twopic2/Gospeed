package main

import (
	"os"
	"sync"
	"time"
)

type writeResults struct {
	Duration  time.Duration
	bytesTime float64
	Error     error
}

type readResults struct {
	Duration  time.Duration
	bytesTime float64
	data      []byte
	Error     error
}

func write(data []byte, key []byte, file string) (time.Duration, float64, error) {
	start := time.Now()

	encryptedData, err := encrypt(data, key)
	checkError(err)

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
	checkError(err)
	defer f.Close()

	_, err = f.Write(encryptedData)
	checkError(err)

	err = f.Sync()
	checkError(err)

	duration := time.Since(start)
	bytesTime := float64(len(encryptedData))

	return duration, bytesTime, nil
}

func writeConcurrent(data []byte, key []byte, file string, wg *sync.WaitGroup, resultChannel chan<- writeResults) {
	defer wg.Done()

	start := time.Now()

	encryptedData, err := encrypt(data, key)
	checkError(err)

	writenfiles := file + "." + string(rune(start.UnixNano()))

	f, err := os.OpenFile(writenfiles, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
	checkError(err)

	defer f.Close()
	defer os.Remove(writenfiles)

	_, err = f.Write(encryptedData)
	checkError(err)

	err = f.Sync()
	checkError(err)

	duration := time.Since(start)
	bytesTime := float64(len(encryptedData))

	resultChannel <- writeResults{
		Duration:  duration,
		bytesTime: bytesTime,
		Error:     nil,
	}
}

func read(file string, key []byte) (time.Duration, float64, []byte, error) {
	start := time.Now()

	f, err := os.OpenFile(file, os.O_RDONLY|os.O_SYNC, 0)
	checkError(err)
	defer f.Close()

	stat, err := f.Stat()
	checkError(err)

	encryptedData := make([]byte, stat.Size())
	_, err = f.Read(encryptedData)
	checkError(err)

	decryptData, err := decrypt(encryptedData, key)
	checkError(err)

	duration := time.Since(start)
	bytesTime := float64(len(encryptedData))

	return duration, bytesTime, decryptData, nil
}

func readConcurrent(file string, key []byte, wg *sync.WaitGroup, resultChannel chan<- readResults) {
	defer wg.Done()

	start := time.Now()

	f, err := os.OpenFile(file, os.O_RDONLY|os.O_SYNC, 0)
	checkError(err)
	defer f.Close()

	stat, err := f.Stat()
	checkError(err)

	encryptedData := make([]byte, stat.Size())
	_, err = f.Read(encryptedData)
	checkError(err)

	decryptedData, err := decrypt(encryptedData, key)
	checkError(err)

	duration := time.Since(start)
	bytesTime := float64(len(encryptedData))

	resultChannel <- readResults{
		Duration:  duration,
		bytesTime: bytesTime,
		data:      decryptedData,
		Error:     nil,
	}
}

func latency(data []byte, key []byte, file string) (time.Duration, error) {
	start := time.Now()

	_, _, err := write(data, key, file)

	checkError(err)

	_, _, _, err = read(file, key)
	checkError(err)

	return time.Since(start), nil
}

/*
func slowWrite(data []byte, key []byte, file string) (time.Duration, float64, error) {
	start := time.Now()

	encryptedData, err := weakEncrypt(data, key)
	checkError(err)

	err = os.WriteFile(file, encryptedData, 0644)
	checkError(err)

	duration := time.Since(start)
	bytesTime := float64(len(encryptedData))

	return duration, bytesTime, nil
}

func slowRead(file string, key []byte) (time.Duration, float64, []byte, error) {
	start := time.Now()

	encryptedData, err := os.ReadFile(file)
	checkError(err)

	decryptData, err := weakDecrypt(encryptedData, key)
	checkError(err)

	duration := time.Since(start)
	bytesTime := float64(len(encryptedData))

	return duration, bytesTime, decryptData, nil
}

func slowLatency(data []byte, key []byte, file string) (time.Duration, err) {
}
*/
