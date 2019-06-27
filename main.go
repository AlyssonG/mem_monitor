package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	Alloc uint64
	Total uint64
	Sys   uint64
	Count uint64
)

func main() {
	Count = 1
	GetMemUsage()

	args := os.Args
	cmd := exec.Command(args[1], args[2:]...)
	cmd.Stdout = os.Stdout

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		for range ticker.C {
			fmt.Println("hey")
			GetMemUsage()
		}

	}()

	cmd.Run()
	cmd.Wait()
	PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func GetMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	Alloc += bToMb(m.Alloc)
	Total += bToMb(m.TotalAlloc)
	Sys += bToMb(m.Sys)
	Count++
}

func PrintMemUsage() {
	fmt.Printf("Alloc = %v MiB", bToMb(Alloc/Count))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(Total/Count))
	fmt.Printf("\tSys = %v MiB", bToMb(Sys/Count))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
