package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var usage = `Usage: portkill <port>`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	port := os.Args[1]
	_, err := strconv.ParseUint(port, 10, 64)
	if err != nil {
		fmt.Println(usage)
		return
	}

	// cmd := exec.Command("netstat", "-ano")
	cmd := exec.Command("cmd", "/c", "netstat -ano | findstr :"+port+" ")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) == 0 {
			fmt.Printf("Port %s not found\n", port)
			return
		}
		fmt.Printf("Error running netstat: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf(":%s ", port)) {
			parts := strings.Fields(line)
			pid := parts[len(parts)-1]
			killCmd := exec.Command("taskkill", "/F", "/PID", pid)
			err = killCmd.Run()
			if err != nil {
				fmt.Printf("Error killing PID %s: %v\n", pid, err)
				return
			}
			fmt.Printf("Killed process with PID %s on port %s\n", pid, port)
			break
		}
	}
}
