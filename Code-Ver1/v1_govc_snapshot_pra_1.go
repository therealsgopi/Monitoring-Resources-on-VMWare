package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/govmomi"
	//"github.com/vmware/govmomi/tree/main/find"
	//"github.com/vmware/govmomi/tree/main/vim25/mo"
    "github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
)

func main() {
	var (
		vcenterURL  = flag.String("url", "", "vCenter URL")
		username    = flag.String("username", "", "vCenter username")
		password    = flag.String("password", "", "vCenter password")
		vmName      = flag.String("vm", "", "Name of the VM to monitor")
		maxSnapshot = flag.Int64("max-snapshot-size", 1073741824, "Maximum snapshot size in bytes")
	)

	flag.Parse()

	if *vcenterURL == "" || *username == "" || *password == "" || *vmName == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u, err := url.Parse(*vcenterURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse vCenter URL: %v\n", err)
		os.Exit(1)
	}

	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create vSphere client: %v\n", err)
		os.Exit(1)
	}

	finder := find.NewFinder(client.Client, true)

	vm, err := finder.VirtualMachine(ctx, *vmName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to find VM %q: %v\n", *vmName, err)
		os.Exit(1)
	}

	var vmMo mo.VirtualMachine
	err = vm.Properties(ctx, vm.Reference(), []string{"snapshot"}, &vmMo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to retrieve VM properties: %v\n", err)
		os.Exit(1)
	}

    // Create a new property collector object
    propColl := view.NewPropertyCollector(client.Client)

    // Define the properties you want to retrieve
    properties := []string{"snapshot.config.memorySizeMB", "snapshot.config.storageSize"}

    // Retrieve the properties for the VM's snapshot
    var s mo.VirtualMachineSnapshot
    err = propColl.RetrieveOne(ctx, ref, properties, &s)
    if err != nil {
        log.Fatalf("Error retrieving snapshot properties: %s", err)
    }

    // Get the snapshot size in bytes
    sizeBytes := s.Snapshot.Config.StorageSize
    

/*	snapshotSize := int64(0)
	for _, s := range vmMo.Snapshot.RootSnapshotList {
		snapshotSize += s.Snapshot.Size
	}*/

	if sizeBytes > *maxSnapshot {
		fmt.Printf("Snapshot size of VM %q exceeds maximum threshold of %d bytes\n", *vmName, *maxSnapshot)
		// Add your alerting code here, such as sending an email or SMS notification
	} 
}

