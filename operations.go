package main

import (
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ncw/directio"
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

var fileCounter int64

func alignDataForDirectIO(data []byte) []byte {
	blockSize := directio.BlockSize
	requiredSize := ((len(data) + blockSize - 1) / blockSize) * blockSize
	alignedData := directio.AlignedBlock(requiredSize)
	copy(alignedData, data)

	for i := len(data); i < requiredSize; i++ {
		alignedData[i] = 0
	}

	return alignedData
}

func write(data []byte, key []byte, file string) (time.Duration, float64, error) {
	start := time.Now()

	encryptedData, err := encrypt(data, key)
	if err != nil {
		return 0, 0, err
	}

	alignedData := alignDataForDirectIO(encryptedData)

	oFile, err := directio.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return 0, 0, err
	}
	defer oFile.Close()

	_, err = oFile.Write(alignedData)
	if err != nil {
		return 0, 0, err
	}

	err = oFile.Sync()
	if err != nil {
		return 0, 0, err
	}

	duration := time.Since(start)
	bytesTime := float64(len(encryptedData))

	return duration, bytesTime, nil
}

func writeConcurrent(data []byte, key []byte, wg *sync.WaitGroup, resultChannel chan<- writeResults) {
	defer wg.Done()

	start := time.Now()

	encryptedData, err := encrypt(data, key)
	if err != nil {
		resultChannel <- writeResults{Error: err}
		return
	}

	counter := atomic.AddInt64(&fileCounter, 1)
	writtenFiles := "tmp" + strconv.FormatInt(counter, 10)

	alignedData := alignDataForDirectIO(encryptedData)

	f, err := directio.OpenFile(writtenFiles, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		resultChannel <- writeResults{Error: err}
		return
	}

	defer func() {
		f.Close()
		os.Remove(writtenFiles)
	}()

	_, err = f.Write(alignedData)
	if err != nil {
		resultChannel <- writeResults{Error: err}
		return
	}

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

	oFile, err := directio.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return 0, 0, nil, err
	}
	defer oFile.Close()

	stat, err := oFile.Stat()
	if err != nil {
		return 0, 0, nil, err
	}

	fileSize := int(stat.Size())
	alignedBuffer := directio.AlignedBlock(fileSize)

	n, err := oFile.Read(alignedBuffer)
	if err != nil {
		return 0, 0, nil, err
	}

	actualSize := n
	for i := n - 1; i >= 0; i-- {
		if alignedBuffer[i] != 0 {
			actualSize = i + 1
			break
		}
	}

	encryptedDatas := alignedBuffer[:actualSize]
	decryptData, err := decrypt(encryptedDatas, key)
	if err != nil {
		decryptData, err = decrypt(alignedBuffer[:n], key)
		if err != nil {
			return 0, 0, nil, err
		}
	}

	duration := time.Since(start)
	bytesTime := float64(len(encryptedDatas))

	return duration, bytesTime, decryptData, nil
}

func readConcurrent(file string, key []byte, wg *sync.WaitGroup, resultChannel chan<- readResults) {
	defer wg.Done()

	start := time.Now()

	oFile, err := directio.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		resultChannel <- readResults{Error: err}
		return
	}
	defer oFile.Close()

	stat, err := oFile.Stat()
	if err != nil {
		resultChannel <- readResults{Error: err}
		return
	}

	fileSize := int(stat.Size())
	alignedBuffer := directio.AlignedBlock(fileSize)

	n, err := oFile.Read(alignedBuffer)
	if err != nil {
		resultChannel <- readResults{Error: err}
		return
	}

	actualSize := n
	for i := n - 1; i >= 0; i-- {
		if alignedBuffer[i] != 0 {
			actualSize = i + 1
			break
		}
	}

	encryptedData := alignedBuffer[:actualSize]

	decryptedData, err := decrypt(encryptedData, key)
	if err != nil {
		decryptedData, err = decrypt(alignedBuffer[:n], key)
		if err != nil {
			resultChannel <- readResults{Error: err}
			return
		}
	}

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
	if err != nil {
		return 0, err
	}

	_, _, _, err = read(file, key)
	if err != nil {
		return 0, err
	}

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
