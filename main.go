package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type DomainDisk struct {
	Device string `xml:"device,attr"`
	Driver struct {
		Name string `xml:"name,attr"`
		Type string `xml:"type,attr"`
	} `xml:"driver"`
	Source struct {
		File string `xml:"file,attr"`
	} `xml:"source"`
	Target struct {
		Dev string `xml:"dev,attr"`
	} `xml:"target"`
}

type Domain struct {
	Devices struct {
		Disks []DomainDisk `xml:"disk"`
	} `xml:"devices"`
}

func main() {
	if !commandExists("qemu-img") || !commandExists("virsh") || !commandExists("sudo") {
		fmt.Println("Required commands (qemu-img, virsh, and/or sudo) are not available. Please install them and try again.")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <VM_NAME>\n", os.Args[0])
		os.Exit(1)
	}

	vmName := os.Args[1]

	if !isVMStopped(vmName) {
		fmt.Printf("Error: VM '%s' is not stopped. Please stop the VM before shrinking its disk.\n", vmName)
		os.Exit(1)
	}

	xmlOutput, err := exec.Command("sudo", "virsh", "dumpxml", vmName).Output()
	if err != nil {
		fmt.Printf("Error getting VM XML: %v\n", err)
		os.Exit(1)
	}

	var domain Domain
	err = xml.Unmarshal(xmlOutput, &domain)
	if err != nil {
		fmt.Printf("Error parsing XML: %v\n", err)
		os.Exit(1)
	}

	for _, disk := range domain.Devices.Disks {
		if disk.Driver.Type == "qcow2" {
			fmt.Printf("Found QCOW2 image: %s\n", disk.Source.File)
			shrinkImage(disk.Source.File)
		}
	}
}

func shrinkImage(imagePath string) {
	fmt.Printf("Shrinking %s...\n", imagePath)

	tempImagePath := imagePath + ".compressed"

	err := exec.Command("sudo", "qemu-img", "convert", "-c", "-f", "qcow2", "-O", "qcow2", imagePath, tempImagePath).Run()
	if err != nil {
		fmt.Printf("Error shrinking image: %v\n", err)
		return
	}

	err = exec.Command("sudo", "mv", tempImagePath, imagePath).Run()
	if err != nil {
		fmt.Printf("Error replacing original image with compressed one: %v\n", err)
		return
	}

	fmt.Printf("Successfully shrank %s\n", imagePath)
}

func isVMStopped(vmName string) bool {
	output, err := exec.Command("sudo", "virsh", "list", "--name", "--state-running").Output()
	if err != nil {
		fmt.Printf("Error checking VM state: %v\n", err)
		return false
	}

	runningVMs := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, vm := range runningVMs {
		if vm == vmName {
			return false
		}
	}
	return true
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
