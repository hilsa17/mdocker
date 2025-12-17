
# mdocker – Minimal Container Runtime
A minimal container runtime, written in Go, that uses core Linux primitives to isolate and run processes. To 
demonstrate how Docker-like containers work internally by creating a lightweight system that launches applications 
inside independent namespaces and a custom root filesystem.
## Features 
- **Namespace Isolation**: Isolate processes using Linux namespaces (PID, UTS, MNT; optional USER/NET)
- **Separate Root Filesystem**: Support running with a separate rootfs using `pivot_root` and populate `/proc`, `/dev`
- **Cgroups Resource Limits**: Set basic resource limits using cgroups (CPU, memory, pids)
- **Re-exec Child Mode**: Parent prepares environment then exec/clone into child mode to build isolated process tree
- **I/O forwarding**: Forward the container’s stdin, stdout, and stderr to the host terminal
## Requirements
- Linux
- BusyBox static
- Go
- cgroups v2 enabled

## Usage
```bash
go build -o mdocker
sudo ./mdocker run /bin/sh



