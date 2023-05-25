package main

import (
	"flag"
	"fmt"
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

var (
	vm string
	action string
	snaps = []Snapshot{}
)

func sizeToMB(size string) float64 {
	sizeMB, _ := strconv.ParseFloat(size[0:len(size)-2], 8)
	if string(size[len(size)-2]) == "K" {
		sizeMB /= 1024
	} else if string(size[len(size)-2]) == "G" {
		sizeMB *= 1024
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

func getVMSnapDetails() {
	//resetting structure "snaps"
	snaps = snaps[0:0]

	// Specify the "govc" command and its arguments
	cmd_ID := exec.Command("govc", "snapshot.tree", "-vm", vm, "-i")
	cmd_name := exec.Command("govc", "snapshot.tree", "-vm", vm)
	cmd_size := exec.Command("govc", "snapshot.tree", "-vm", vm, "-s")
	cmd_crDate := exec.Command("govc", "snapshot.tree", "-vm", vm, "-D")

	// Run the commands and capture the output
	output_ID, _ := cmd_ID.Output()
	output_name, _ := cmd_name.Output()
	output_size, _ := cmd_size.Output()
	output_crDate, _ := cmd_crDate.Output()

	lines_ID := strings.Split(strings.TrimSuffix(string(output_ID), "\n"), "\n")
	lines_name := strings.Split(strings.TrimSuffix(string(output_name), "\n"), "\n")
	lines_size := strings.Split(strings.TrimSuffix(string(output_size), "\n"), "\n")
	lines_crDate := strings.Split(strings.TrimSuffix(string(output_crDate), "\n"), "\n")

	snaps = make([]Snapshot, len(lines_ID))

	storeSnapDetails(lines_ID, "ID")
	storeSnapDetails(lines_name, "name")
	storeSnapDetails(lines_size, "size")
	storeSnapDetails(lines_crDate, "crDate")
}

func storeSnapDetails(lines []string, detail string) {
	for i := 0; i < len(lines); i++ {
		if detail == "name" {
			snaps[i].name = strings.TrimSpace(lines[i])
		} else {
			startInd := strings.Index(lines[i], "[") + 1
			endInd := strings.Index(lines[i], "]")
			value := lines[i][startInd:endInd]
			if detail == "ID" {
				snaps[i].id = value
			}
			if detail == "size" {
				snaps[i].size = sizeToMB(value)
			}
			if detail == "crDate" {
				layout := "Jan 1 15:04"
				crDate, _ := time.Parse(layout, value)
				snaps[i].date = crDate
			}
		}
	}
}

func dispSnapDetails() {
	fmt.Println("Number of snapshots: ", len(snaps))
	for snap := range snaps {
		fmt.Println("Snapshot", snap, ":-")
		fmt.Println(snaps[snap].id)
		fmt.Println(snaps[snap].size)
		fmt.Println(snaps[snap].date)
		fmt.Println(snaps[snap].name)
		fmt.Println()
	}
}

func deleteSnapFromStruct(index int) {
	snaps = append(snaps[:index], snaps[index+1:]...)
}

func checkSnapshots(action string) {
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
				output_snap_rem, _ := cmd_snap_rem.Output()
				deleteSnapFromStruct(snap)
				fmt.Printf("ALERT: Snapshot %v of VM %v successfully deleted %v\n", snaps[snap].name, vm, output_snap_rem)
			} else {
				fmt.Printf("WARNING: Snapshot %v of VM %v will automatically be deleted after 5 days\n", snaps[snap].name, vm)
			}
		}
	}
}

func init() {
	flag.StringVar(&vm, "vm", "Host2_Mint3", "name of the vm whose snapshots are to be checked")
	flag.StringVar(&action, "action", "warn", "specify if snapshots are to be deleted or warned")
	flag.Parse()
}

func main() {
	getVMSnapDetails()
	
	fmt.Println("Details of Snapshots of VM", vm, " before Checking:-")	
	dispSnapDetails()

	checkSnapshots(action)
	
	if action == "delete" {
		fmt.Println("Details of Snapshots of VM", vm, " after Checking:-")	
		dispSnapDetails()
	}
}
