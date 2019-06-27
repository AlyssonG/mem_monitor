package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	MemUsed uint64
	Count   uint64
)

func main() {
	Count = 1

	args := os.Args
	if len(args) < 4 {
		fmt.Println("There are no parameters enough to run monitor")
		os.Exit(1)
	}

	mili, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Unexpected value for monitor time:", args[1])
		os.Exit(1)
	}

	cmd := exec.Command(args[2], args[3:]...)
	cmd.Stdout = os.Stdout

	go func() {
		ticker := time.NewTicker(time.Duration(mili) * time.Millisecond)
		for range ticker.C {
			GetMemUsage(cmd.Process.Pid)
		}

	}()

	cmd.Run()
	cmd.Wait()
	PrintMemUsage()
}

func GetMemUsage(pid int) {
	mem, _ := calculateMemory(pid)
	MemUsed += mem
	Count++
}

func PrintMemUsage() {
	fmt.Printf("MemUsed = %v KB\n", MemUsed/Count)
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
