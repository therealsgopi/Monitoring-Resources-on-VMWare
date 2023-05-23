package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// Specify the "govc" command and its arguments
	// cmd := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3")
	cmd := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-D", "-s", "-i")

	// Run the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to run govc command: %v", err)
	}

	// Convert the output to a string and print it
	fmt.Println(string(output))
}
