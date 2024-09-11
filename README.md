# kvm-qcow-shrink

`kvm-qcow-shrink` is a command-line tool written in Go that compresses QCOW2 image files used by KVM virtual machines. It automatically identifies QCOW2 images associated with a specified VM and efficiently compresses them.

## Features

- Automatically identifies QCOW2 images for a specified KVM virtual machine
- Verifies that the VM is stopped before proceeding
- Uses an efficient single-step compression process
- Compresses QCOW2 images in place to optimize disk space
- Uses sudo to execute commands with necessary permissions

## Notes

- This tool requires sudo privileges to run. Make sure you have the necessary permissions.
- Ensure the VM is stopped before running this tool. The tool will check this and exit if the VM is running.
- Always ensure you have a backup of your VM and its disk images before using this tool.
- The compression process may take some time depending on the size of your QCOW2 images.
- This tool modifies the original image files in place. Make sure you have enough free space on your disk to accommodate temporary files during the compression process.

## Prerequisites

Before using kvm-qcow-shrink, ensure you have the following installed on your system:

- `virsh` command-line tool
- `qemu-img` command-line tool
- `sudo` access

## Installation

Build the tool:

```
$ go build -o kvm-qcow-shrink
```

## Usage

To use kvm-qcow-shrink, run the following command with sudo:

```
$ kvm-qcow-shrink <VM_NAME>
```

Replace `<VM_NAME>` with the name of the KVM virtual machine whose images you want to compress.

Example:
```
$ kvm-qcow-shrink ubuntu-server
```

The tool will then:
1. Check if the specified VM is stopped
2. Retrieve the VM's XML configuration
3. Identify all QCOW2 image files associated with the VM
4. Compress each QCOW2 image in place

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
