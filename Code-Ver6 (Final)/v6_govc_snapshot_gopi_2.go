package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Snapshot struct {
	name string
	id   string
	size float64
	date time.Time
}

func sizeToMB(size string) float64 {
	var sizeMB float64
	if string(size[len(size)-2]) == "K" {
		sizeMB, _ = strconv.ParseFloat(size[0:len(size)-2], 8)
		sizeMB /= 1024
	} else if string(size[len(size)-2]) == "G" {
		sizeMB, _ = strconv.ParseFloat(size[0:len(size)-2], 8)
		sizeMB *= 1024
	} else {
		sizeMB, _ = strconv.ParseFloat(size[0:len(size)-2], 8)
	}
	return sizeMB
}

func currentDate() time.Time {
	date := time.Now()
	loc, _ := time.LoadLocation("UTC")
	date = date.In(loc)
	date = time.Date(0000, date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, date.Location())
	return date
}

func snapLife(creationDate time.Time) int64 {
	snap_lifespan := currentDate().Sub(creationDate)
	days := int64(snap_lifespan.Hours() / 24)
	if days < 0 {
		days = 365 + days
	}
	return days
}

func into_struct(snaps []Snapshot, lines []string, d int) {
	layout := "Jan 2 15:04"
	for i := 0; i < len(lines); i++ {
		startInd := strings.Index(lines[i], "[") + 1
		endInd := strings.Index(lines[i], "]")
		value := lines[i][startInd:endInd]
		if d == 0 {
			snaps[i].id = value
		}
		if d == 1 {
			snaps[i].size = sizeToMB(value)
		}
		if d == 2 {
			t, err := time.Parse(layout, value)
			if err != nil {
				fmt.Println("Error parsing input string:", err)
				return
			}
			snaps[i].date = t
		}
	}
}

func checkSnapshots(snaps []Snapshot, action string) {
	for snap := range snaps {
		var snap_rem_flag int64
		snapDays := snapLife(snaps[snap].date)
		if snaps[snap].size > 5120 {
			if snapDays > 3 {
				snap_rem_flag = 1
			}
		} else if snaps[snap].size > 1024 {
			if snapDays > 30 {
				snap_rem_flag = 1
			}
		} else {
			if snapDays > 180 {
				snap_rem_flag = 1
			}
		}

		if snap_rem_flag == 1 {
			if action == "delete" {
				cmd_snap_rem := exec.Command("govc", "snapshot.remove", "-vm", vm, snaps[snap].name)
				output_snap_rem, err_snap_rem := cmd_snap_rem.Output()
				if err_snap_rem != nil {
					log.Fatalf("Failed to run govc command: %v", err_snap_rem)
				}
				fmt.Printf("ALERT: Snapshot %v of VM %v successfully deleted %v\n", snaps[snap].name, vm, output_snap_rem)
			} else {
				fmt.Printf("WARNING: Snapshot %v of VM %v will automatically be deleted after 5 days\n", snaps[snap].name, vm)
			}
		}
	}
}

var (
	vm string
	action string
)

func init() {
	flag.StringVar(&vm, "vm", "Host2_Mint3", "name of the vm whose snapshots are to be checked")
	flag.StringVar(&action, "action", "warn", "specify if snapshots are to be deleted or warned")
	flag.Parse()
}

func main() {
	// Specify the "govc" command and its arguments
	cmd_d := exec.Command("govc", "snapshot.tree", "-vm", vm, "-D")
	cmd_i := exec.Command("govc", "snapshot.tree", "-vm", vm, "-i")
	cmd_s := exec.Command("govc", "snapshot.tree", "-vm", vm, "-s")
	cmd_n := exec.Command("govc", "snapshot.tree", "-vm", vm)

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

	lines_i := strings.Split(strings.TrimSuffix(string(output_i), "\n"), "\n")
	lines_s := strings.Split(strings.TrimSuffix(string(output_s), "\n"), "\n")
	lines_d := strings.Split(strings.TrimSuffix(string(output_d), "\n"), "\n")
	lines_n := strings.Split(strings.TrimSuffix(string(output_n), "\n"), "\n")

	snaps := make([]Snapshot, len(lines_n))

	for i, line := range lines_n {
		snaps[i].name = strings.TrimSpace(line)
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

	/*
	       cmd_snap_rem := exec.Command("govc", "snapshot.remove", "-vm", "Host2_Mint3", "trial1_1")
	       output_snap_rem, err_snap_rem := cmd_snap_rem.Output()
	       if err_snap_rem != nil {
	   		log.Fatalf("Failed to run govc command: %v", err_snap_rem)
	   	}
	   	fmt.Println(output_snap_rem)
	*/
	checkSnapshots(snaps, action)
	fmt.Println(vm, action)

}


