package main

import (
	"fmt"
	"strings"

	"log"
	"os/exec"
	"time"
)

type Snapshot struct {
	name string
	id   string
	size string
	date time.Time
}

func main() {
	// Specify the "govc" command and its arguments
	cmd_d := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-D")
	    cmd_i := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-i")
	    cmd_s := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-s")
        cmd_n := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3")

		// Run the command and capture the output
		output_d, err_d := cmd_d.Output()
	    output_i, err_i := cmd_i.Output()
	    output_s, err_s := cmd_s.Output()
        output_n, err_s := cmd_n.Output()
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
	snaps := make([]Snapshot, 5)
	/*output_i := "[snapshot-41]  trial1\n [snapshot-42]  trial1_1\n [snapshot-43]  trial2"
	output_s := "[19.6KB]  trial1\n [1.0MB]  trial1_1\n [169.3MB]  trial2"
	output_d := "[May 8 11:38]  trial1\n [May 8 11:45]  trial1_1\n [May 8 12:02]  trial2"*/

	lines_i := strings.Split(strings.TrimSuffix(string(output_i),"\n"), "\n")
	for _, line := range lines_i {
		fmt.Println(line)
	}
	lines_s := strings.Split(strings.TrimSuffix(string(output_s),"\n"), "\n")
    for _, line := range lines_s {
		fmt.Println(line)
	}
	lines_d := strings.Split(strings.TrimSuffix(string(output_d),"\n"), "\n")
    for _, line := range lines_d {
		fmt.Println(line)
	}
    lines_n := strings.Split(strings.TrimSuffix(string(output_n),"\n"), "\n")
    for i, line := range lines_n {
		fmt.Println(line)
        snaps[i].name=line
	}


	temp := ""
	for i := 0; i < len(lines_i); i++ {
		// Loop over each character in the string
		for j := 0; j < len(lines_i[i]); j++ {
			if lines_i[i][j] == '[' {
				j++
				for lines_i[i][j] != ']' {
					//fmt.Println(string(lines_i[i][j]))
					temp += string(lines_i[i][j])
					j++
				}
			}
		}
		snaps[i].id = temp
		temp = ""
	}


	for i := 0; i < len(lines_s); i++ {
		// Loop over each character in the string
		for j := 0; j < len(lines_s[i]); j++ {
			if lines_s[i][j] == '[' {
				j++
				for lines_s[i][j] != ']' {
					//fmt.Println(string(lines_i[i][j]))
					temp += string(lines_s[i][j])
					j++
				}
			}
		}
		snaps[i].size = temp
		temp = ""
	}

	layout := "Jan 2 15:04"
	for i := 0; i < len(lines_d); i++ {
		// Loop over each character in the string
		for j := 0; j < len(lines_d[i]); j++ {
			if lines_d[i][j] == '[' {
				j++
				for lines_d[i][j] != ']' {
					//fmt.Print(string(lines_d[i][j]))
					temp += string(lines_d[i][j])
					j++
				}
			}
		}
		t, err := time.Parse(layout, temp)
		if err != nil {
			fmt.Println("Error parsing input string:", err)
			return
		}
		snaps[i].date = t
		temp = ""
	}

	for snap := range snaps {
		fmt.Println(snaps[snap].id)
		fmt.Println(snaps[snap].size)
		fmt.Println(snaps[snap].date)
	}
}
