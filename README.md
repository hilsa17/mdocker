
# mdocker – Minimal Container Runtime
A minimal container runtime, written in Go, that uses core Linux primitives to isolate and run processes. To 
demonstrate how Docker-like containers work internally by creating a lightweight system that launches applications 
inside independent namespaces and a custom root filesystem.
## Features 
- **Namespace Isolation**: Isolate processes using Linux namespaces (PID, UTS, MNT; optional USER/NET). This was implemented using:
  `syscall.CLONE_NEWPID
   syscall.CLONE_NEWUTS
   syscall.CLONE_NEWNS
   syscall.CLONE_NEWNET`

- **Separate Root Filesystem**: Support running with a separate rootfs using `pivot_root` and populate `/proc`, `/dev`
- **Cgroups Resource Limits**: Set basic resource limits using cgroups (CPU, memory, pids).
  A dedicated cgroup is created at `/sys/fs/cgroup/mdocker` which is configured via `cpu.max`, `memory.max`, `pids.max`, `cgroup.procs`. This can be verified:
  `cat /sys/fs/cgroup/mdocker/pids.max`
  `cat /sys/fs/cgroup/mdocker/pids.current`
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
```
This project was developed on WSL2. Due to WSL2 limitations, container networking couldn't be implemented but network namespace was created, which can be verified using `ip addr`.
## Setup Instructions
1. Clone Repository
```bash
git clone https://github.com/hilsa17/mdocker
cd mdocker
```
2. Prepare Root Filesystem
```bash
mkdir -p rootfs/{bin,proc,sys,dev,tmp}
cp /bin/busybox rootfs/bin/
```

Create BusyBox symlinks:
```bash
cd rootfs/bin
./busybox --install .
cd ../../
```

Ensure permissions:
```bash
chmod -R 755 rootfs
```
3. Build mdocker
```bash
go build -o mdocker
```
5. Run Container
```bash
sudo ./mdocker run /bin/sh
```
##Demo commands
1. PID isolation
```bash
ps
```
2. Process limit demo
```bash
for i in $(seq 1 30); do sh -c "sleep 100" & done
```

3. Check:
```bash
cat /sys/fs/cgroup/mdocker/pids.current
```

