package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

type Result struct {
	Size            int
	WriteTime       time.Duration
	ReadTime        time.Duration
	LatencyTime     time.Duration
	WriteThroughput float64
	ReadThroughput  float64
}

type megaByte int

const (
	megaByteF       float64  = 1024 * 1024
	regular         megaByte = 1024 * 1024
	tenMegabyte     megaByte = 10 * 1024 * 1024
	hundredMegabyte megaByte = 100 * 1024 * 1024
	gigaByte        megaByte = 1000 * 1024 * 1024
)

func runAES(sizes []int, filename string) []Result {
	numCpus := runtime.NumCPU()

	var mu sync.Mutex

	result := make([]Result, len(sizes))
	key := make([]byte, 32) // 256bit key
	rand.Read(key)

	for i, size := range sizes {
		data := make([]byte, size)
		rand.Read(data)

		_, _, err := write(data, key, filename)
		if err != nil {
			fmt.Printf("Error writing initial file: %v\n", err)
			continue
		}

		numWorkers := min(numCpus, len(data))
		coreChunks := spiltCores(data, numWorkers)

		var wg sync.WaitGroup

		// Go routines make
		writeResultsChannel := make(chan writeResults, numCpus)

		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(i int, chunk []byte) {
				mu.Lock()
				writeConcurrent(chunk, key, filename, &wg, writeResultsChannel)
				mu.Unlock()
			}(w, coreChunks[w])
		}
		wg.Wait()

		totalWriteDuration := time.Duration(0)
		for i := 0; i < numWorkers; i++ {
			writeResult := <-writeResultsChannel
			totalWriteDuration += writeResult.Duration
		}

		readResultsChannel := make(chan readResults, numCpus)
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(i int, filename string) {
				mu.Lock()
				readConcurrent(filename, key, &wg, readResultsChannel)
				mu.Unlock()
			}(w, filename)
		}
		wg.Wait()

		totalReadDurration := time.Duration(0)
		for i := 0; i < numWorkers; i++ {

			readResult := <-readResultsChannel
			totalReadDurration += readResult.Duration
		}

		latencyTime, _ := latency(data, key, filename)

		writeThroughput := float64(size) / totalWriteDuration.Seconds() / megaByteF
		readThroughput := float64(size) / totalReadDurration.Seconds() / megaByteF

		result[i] = Result{
			Size:            size,
			WriteTime:       totalWriteDuration,
			ReadTime:        totalReadDurration,
			LatencyTime:     latencyTime,
			WriteThroughput: writeThroughput,
			ReadThroughput:  readThroughput,
		}
	}
	return result
}

func testAES() {
	fmt.Println("Welcome to Gopeed! A basic encryption file transfer benchmark using AES.\nBeware Evil femboys are stealing your data! Your data must be encrypted properly! With uwu-AES-uwu everything will be saved!")
	file := "encryption_test.txt"
	defer os.Remove(file)
	sizes := []int{
		int(regular),
		int(tenMegabyte),
		int(hundredMegabyte),
		int(gigaByte),
	}
	result := runAES(sizes, file)

	fmt.Println("Size (bytes) | Write (MB/s) | Read (MB/s) | Latency (ms)")
	for _, r := range result {
		fmt.Printf("%-13d | %-12.2f | %-11.2f | %-12.2f\n",
			r.Size,
			r.WriteThroughput,
			r.ReadThroughput,
			float64(r.LatencyTime.Microseconds())/1000.0)
	}
	fmt.Print("Quickly we discord kittens must report to our masters for our discord calls!")

	art := `
⢸⠂⠀⠀⠀⠘⣧⠀⠀⣟⠛⠲⢤⡀⠀⠀⣰⠏⠀⠀⠀⠀⠀⢹⡀
⠀⡿⠀⠀⠀⠀⠀⠈⢷⡀⢻⡀⠀⠀⠙⢦⣰⠏⠀⠀⠀⠀⠀⠀⢸⠀
⠀⡇⠀⠀⠀⠀⠀⠀⢀⣻⠞⠛⠀⠀⠀⠀⠻⠀⠀⠀⠀⠀⠀⠀⢸⠀
⠀⡇⠀⠀⠀⠀⠀⠀⠛⠓⠒⠓⠓⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⠀
⠀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣸⠀
⠀⢿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⠀⠀⢀⡟⠀
⠀⠘⣇⠀⠘⣿⠋⢹⠛⣿⡇⠀⠀⠀⠀⣿⣿⡇⠀⢳⠉⠀⣠⡾⠁⠀
⣦⣤⣽⣆⢀⡇⠀⢸⡇⣾⡇⠀⠀⠀⠀⣿⣿⡷⠀⢸⡇⠐⠛⠛⣿⠀
⠹⣦⠀⠀⠸⡇⠀⠸⣿⡿⠁⢀⡀⠀⠀⠿⠿⠃⠀⢸⠇⠀⢀⡾⠁⠀
⠀⠈⡿⢠⢶⣡⡄⠀⠀⠀⠀⠉⠁⠀⠀⠀⠀⠀⣴⣧⠆⠀⢻⡄⠀⠀
⠀⢸⠃⠀⠘⠉⠀⠀⠀⠠⣄⡴⠲⠶⠴⠃⠀⠀⠀⠉⡀⠀⠀⢻⡄⠀
⠀⠘⠒⠒⠻⢦⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⠞⠛⠒⠛⠋⠁⠀
⠀⠀⠀⠀⠀⠀⠸⣟⠓⠒⠂⠀⠀⠀⠀⠀⠈⢷⡀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠙⣦⠀⠀⠀⠀⠀⠀⠀⠀⠈⢷⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣼⣃⡀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣆⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠉⣹⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⢻⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⡆`
	fmt.Println(art)
}
