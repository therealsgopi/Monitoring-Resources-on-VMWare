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

func print_str(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func into_struct(snaps []Snapshot, lines []string, d int) {
	layout := "Jan 2 15:04"
	temp := ""
	for i := 0; i < len(lines); i++ {
		// Loop over each character in the string
		for j := 0; j < len(lines[i]); j++ {
			if lines[i][j] == '[' {
				j++
				for lines[i][j] != ']' {
					//fmt.Println(string(lines_i[i][j]))
					temp += string(lines[i][j])
					j++
				}
			}
		}
		if d == 0 {
			snaps[i].id = temp
		}
		if d == 1 {
			snaps[i].size = temp
		}
		if d == 2 {
			t, err := time.Parse(layout, temp)
			if err != nil {
				fmt.Println("Error parsing input string:", err)
				return
			}
			snaps[i].date = t
		}
		temp = ""
	}
}

func main() {
	// Specify the "govc" command and its arguments
	snaps := make([]Snapshot, 6)
	cmd_d := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-D")
	cmd_i := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-i")
	cmd_s := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3", "-s")
	cmd_n := exec.Command("govc", "snapshot.tree", "-vm", "Host2_Mint3")

	// Run the command and capture the output
	output_d, err_d := cmd_d.Output()
	output_i, err_i := cmd_i.Output()
	output_s, err_s := cmd_s.Output()
	output_n, err_n := cmd_n.Output()
	if err_d != nil {
		log.Fatalf("Failed to run govc command: %v", err_d)
	}
	if err_i != nil {
		log.Fatalf("Failed to run govc command: %v", err_i)
	}
	if err_s != nil {
		log.Fatalf("Failed to run govc command: %v", err_s)
	}
	if err_n != nil {
		log.Fatalf("Failed to run govc command: %v", err_n)
	}

	//To try at home uncomment below and comment above
	/*snaps := make([]Snapshot, 3)
	output_i := "[snapshot-41]  trial1\n [snapshot-42]  trial1_1\n [snapshot-43]  trial2"
	output_s := "[19.6KB]  trial1\n [1.0MB]  trial1_1\n [169.3MB]  trial2"
	output_d := "[May 8 11:38]  trial1\n [May 8 11:45]  trial1_1\n [May 8 12:02]  trial2"
	output_n := "trial1\n trial1_1\n trial2"*/

	lines_i := strings.Split(strings.TrimSuffix(string(output_i), "\n"), "\n")
	lines_s := strings.Split(strings.TrimSuffix(string(output_s), "\n"), "\n")
	lines_d := strings.Split(strings.TrimSuffix(string(output_d), "\n"), "\n")
	lines_n := strings.Split(strings.TrimSuffix(string(output_n), "\n"), "\n")
	for i, line := range lines_n {
		snaps[i].name = line
	}

	into_struct(snaps, lines_i, 0)
	into_struct(snaps, lines_s, 1)
	into_struct(snaps, lines_d, 2)

	for snap := range snaps {
		fmt.Println(snaps[snap].id)
		fmt.Println(snaps[snap].size)
		fmt.Println(snaps[snap].date)
		fmt.Println(snaps[snap].name)
	}
}
