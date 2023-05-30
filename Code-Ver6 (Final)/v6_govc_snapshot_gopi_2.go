package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// structure to store individual snapshot details
type Snapshot struct {
	name string
	id   string
	size float64
	date time.Time
}

// declaring and instantiating global variables
var (
	vm string
	action string
	snaps = []Snapshot{}
	illegalSnaps = []int{}
)

// function for converting size of snapshots to MB from KB & GB
func sizeToMB(size string) float64 {
	sizeMB, _ := strconv.ParseFloat(size[0:len(size)-2], 8)
	if string(size[len(size)-2]) == "K" {
		sizeMB /= 1024
	} else if string(size[len(size)-2]) == "G" {
		sizeMB *= 1024
	}
	return sizeMB
}

// function to return current date in the same format as creation dates of the snapshots
func currentDate() time.Time {
	date := time.Now()
	loc, _ := time.LoadLocation("UTC")
	date = date.In(loc)
	date = time.Date(0000, date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, date.Location())
	return date
}

// function to compute the age of a snapshot from current date
func snapLife(creationDate time.Time) int64 {
	snap_lifespan := currentDate().Sub(creationDate)
	days := int64(snap_lifespan.Hours() / 24)
	if days < 0 {
		days = 365 + days
	}
	return days
}

// function to get the snapshot details of a specific VM
func getVMSnapDetails() {
	/*
	output_ID := "[snapshot-41]  trial1\n [snapshot-42]  trial1_1\n [snapshot-43]  trial2"
	output_size := "[5555MB]  trial1\n [550MB]  trial1_1\n [1.25GB]  trial2"
	output_crDate := "[May 25 17:38]  trial1\n [Dec 1 11:45]  trial1_1\n [Apr 15 12:02]  trial2"
	output_name := "trial1\n trial1_1\n trial2"
	*/
	//resetting array SNAPS
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
	
	// splitting the output to form an array of individual lines of details for each snapshot
	lines_ID := strings.Split(strings.TrimSuffix(string(output_ID), "\n"), "\n")
	lines_name := strings.Split(strings.TrimSuffix(string(output_name), "\n"), "\n")
	lines_size := strings.Split(strings.TrimSuffix(string(output_size), "\n"), "\n")
	lines_crDate := strings.Split(strings.TrimSuffix(string(output_crDate), "\n"), "\n")

	// initializing empty SNAPS structure to avoid IndexOutOfBound error
	snaps = make([]Snapshot, len(lines_ID))

	// storing all the read details of snapshots
	storeSnapDetails(lines_ID, "ID")
	storeSnapDetails(lines_name, "name")
	storeSnapDetails(lines_size, "size")
	storeSnapDetails(lines_crDate, "crDate")
}

// function for storing all the details of 
//all the snapshots in the structure SNAPS
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
				layout := "Jan 2 15:04"
				crDate, _ := time.Parse(layout, value)
				snaps[i].date = crDate
			}
		}
	}
}

// function for displaying all the details 
// of all the snapshots in the structure SNAPS
func dispSnapDetails() {
	fmt.Println("Number of snapshots:", len(snaps))
	for snap := range snaps {
		fmt.Println("Snapshot", (snap + 1), ":-")
		fmt.Println("ID:", snaps[snap].id)
		fmt.Println("Name:", snaps[snap].name)
		fmt.Println("Size:", snaps[snap].size)
		fmt.Println("Date:", snaps[snap].date)
		fmt.Println()
	}
}

// function for deleting the snapshot details from the 
// structure SNAPS when it is deleted from the storage
func deleteSnapFromStruct(index int) {
	snaps = append(snaps[:index], snaps[index+1:]...)
}

// function for checking all the snapshots stored in 
// the structure SNAPS and mark the illegal ones
func checkSnapshots() {
	//resetting array illegalSnaps
	illegalSnaps = illegalSnaps[0:0]

	for snap := range snaps {
		snapDays := snapLife(snaps[snap].date)
		if snaps[snap].size > 5120 {
			if snapDays > 3 {
				illegalSnaps = append(illegalSnaps, snap)
			}
		} else if snaps[snap].size > 1024 {
			if snapDays > 30 {
				illegalSnaps = append(illegalSnaps, snap)
			}
		} else {
			if snapDays > 180 {
				illegalSnaps = append(illegalSnaps, snap)
			}
		}
	}
}

// funtion for taking action on the snapshots that are marked illegal
func takeAction() {
	for ind, illegalSnapInd := range illegalSnaps {
		if action == "delete" {
			cmd_snap_rem := exec.Command("govc", "snapshot.remove", "-vm", vm, snaps[illegalSnapInd - ind].name)
			output_snap_rem, _ := cmd_snap_rem.Output()
			fmt.Printf("ALERT: Snapshot %v of VM %v successfully deleted %v\n", snaps[illegalSnapInd - ind].name, vm, output_snap_rem)
			deleteSnapFromStruct(illegalSnapInd - ind)
		} else {
			fmt.Printf("WARNING: Snapshot %v of VM %v will automatically be deleted after 5 days\n", snaps[illegalSnapInd].name, vm)
		}
	}
}

// initializing the values entered for the variables through CLI
func init() {
	flag.StringVar(&vm, "vm", "Host2_Mint3", "name of the vm whose snapshots are to be checked")
	flag.StringVar(&action, "action", "warn", "specify if snapshots are to be deleted or warned")
	flag.Parse()
}

// Driver Code
func main() {
	getVMSnapDetails()
	
	fmt.Println("\nDetails of Snapshots of VM", vm, "before Checking:-")	
	dispSnapDetails()

	checkSnapshots()
	takeAction()
	
	if action == "delete" {
		fmt.Println("\nDetails of Snapshots of VM", vm, "after Checking:-")	
		dispSnapDetails()
	}
}
