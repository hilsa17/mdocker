package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: mdocker run <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
}
func setupCgroup(pid int) {
	cgroupPath := "/sys/fs/cgroup/mdocker"

	// create cgroup
	if err := os.Mkdir(cgroupPath, 0755); err != nil && !os.IsExist(err) {
		panic(err)
	}

	// CPU: 50%
	os.WriteFile(cgroupPath+"/cpu.max", []byte("50000 100000"), 0644)

	// Memory: 100mb
	os.WriteFile(cgroupPath+"/memory.max", []byte("104857600"), 0644)

	// Max 20 processes
	os.WriteFile(cgroupPath+"/pids.max", []byte("20"), 0644)

	// add process to cgroup
	os.WriteFile(cgroupPath+"/cgroup.procs", []byte(fmt.Sprintf("%d", pid)), 0644)
}

func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS|
                        syscall.CLONE_NEWNET,
	}
        if err := cmd.Start(); err != nil {
		panic(err)
	}

	setupCgroup(cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
}

func child() {
	fmt.Println("Running container...")

	// hostname
	if err := syscall.Sethostname([]byte("mdocker")); err != nil {
		panic(err)
	}

	// make mounts private
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		panic(err)
	}

	rootfs, err := filepath.Abs("rootfs")
        if err != nil {
	panic(err)
        }


	// bind mount rootfs
	if err := syscall.Mount(rootfs, rootfs, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		panic(err)
	}

	// create pivot_root directories
	putOld := filepath.Join(rootfs, ".oldroot")
	if err := os.MkdirAll(putOld, 0700); err != nil {
		panic(err)
	}

	// pivot_root
	if err := syscall.PivotRoot(rootfs, putOld); err != nil {
		panic(err)
	}

	if err := os.Chdir("/"); err != nil {
		panic(err)
	}

	// unmount old root
	putOld = "/.oldroot"
	if err := syscall.Unmount(putOld, syscall.MNT_DETACH); err != nil {
		panic(err)
	}
	if err := os.RemoveAll(putOld); err != nil {
		panic(err)
	}

	// mount proc
	if err := os.MkdirAll("/proc", 0555); err != nil {
		panic(err)
	}
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		panic(err)
	}

        // mount cgroups v2
        os.MkdirAll("/sys/fs/cgroup", 0755)
        if err := syscall.Mount("cgroup2", "/sys/fs/cgroup", "cgroup2", 0, ""); err != nil {
                panic(err)
        }

	// run command
	cmd := exec.Command(os.Args[2])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

