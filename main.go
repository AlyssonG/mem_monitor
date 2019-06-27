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
	Alloc += bToKb(m.Alloc)
	Total += bToKb(m.TotalAlloc)
	Sys += bToKb(m.Sys)
	Count++
}

func PrintMemUsage() {
	fmt.Printf("Alloc = %v KiB", bToKb(Alloc/Count))
	fmt.Printf("\tTotalAlloc = %v KiB", bToKb(Total/Count))
	fmt.Printf("\tSys = %v KiB\n", bToKb(Sys/Count))
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
