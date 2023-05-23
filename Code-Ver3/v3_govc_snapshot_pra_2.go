package main

import (
	"fmt"
	"log"
	"os/exec"
    "time"
)

type Snapshot struct {
    name string
    id  string
    size string
    date time.Time
}

func main() {
	// Specify the "govc" command and its arguments
	cmd_d := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-D")
    cmd_i := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-i")
    cmd_s := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-s")

	// Run the command and capture the output
	output_d, err_d := cmd_d.Output()
    output_i, err_i := cmd_i.Output()
    output_s, err_s := cmd_s.Output()
	if err_d != nil {
		log.Fatalf("Failed to run govc command: %v", err_d)
	}
	// Convert the output to a string and print it
	fmt.Println(string(output_d))

    if err_i != nil {
		log.Fatalf("Failed to run govc command: %v", err_i)
	}
	// Convert the output to a string and print it
	fmt.Println(string(output_i))

    if err_s != nil {
		log.Fatalf("Failed to run govc command: %v", err_s)
	}
	// Convert the output to a string and print it
	fmt.Println(string(output_s))
    
    var snaps []Snapshot
    var ct=0
    for i:= range(output_i) {
        if string(output_i[i])=="["{
            for string(output[i])!="]"{
                snaps[ct].id=output_i[i+1]
            } 
        }
        fmt.Print(snaps)
    }
}
