package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var (
	Alloc   uint64
	Total   uint64
	Sys     uint64
	MemUsed uint64
	Count   uint64
)

func main() {
	Count = 1

	args := os.Args
	cmd := exec.Command(args[1], args[2:]...)
	cmd.Stdout = os.Stdout

	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		for range ticker.C {
			GetMemUsage(cmd.Process.Pid)
		}

	}()

	cmd.Run()
	cmd.Wait()
	PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func GetMemUsage(pid int) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	Alloc += bToKb(m.Alloc)
	Total += bToKb(m.TotalAlloc)
	Sys += bToKb(m.Sys)
	MemUsed, _ = calculateMemory(pid)
	Count++
}

func PrintMemUsage() {
	// fmt.Printf("Alloc = %v KiB", bToKb(Alloc/Count))
	fmt.Printf("MemUsed = %v KB\n", MemUsed/Count)
	// fmt.Printf("\tTotalAlloc = %v KiB", bToKb(Total/Count))
	// fmt.Printf("\tSys = %v KiB\n", bToKb(Sys/Count))
}

func calculateMemory(pid int) (uint64, error) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/smaps", pid))
	if err != nil {
		return 0, err
	}
	defer f.Close()

	res := uint64(0)
	pfx := []byte("Pss:")
	r := bufio.NewScanner(f)
	for r.Scan() {
		line := r.Bytes()
		if bytes.HasPrefix(line, pfx) {
			var size uint64
			_, err := fmt.Sscanf(string(line[4:]), "%d", &size)
			if err != nil {
				return 0, err
			}
			res += size
		}
	}
	if err := r.Err(); err != nil {
		return 0, err
	}

	return res, nil
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
