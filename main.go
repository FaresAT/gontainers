package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("invalid command")
	}
}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// setting in/out/err to OS in/out/err to get output
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// CLONE_NEWUTS is what clones the hostname
	// CLONE_NEWPID clones the process ids
	// CLONE_NEWNS generates a new namespace
	// CLONE_NEWUSER creates a new user namespace using the UID mappings
	// UidMappings call allows the container to be run without being a root user
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{{
			ContainerID: 0,
			HostID:      1000,
			Size:        1,
		}},
		Unshareflags: syscall.CLONE_NEWNS,
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	controlGroups()

	// setting new hostname
	syscall.Sethostname([]byte(LoadEnv("CONTAINER_NAME")))

	// setting new root directory, must be some form of linux filesystem (specified in .env file)
	syscall.Chroot(LoadEnv("CONTAINER_ROOT"))
	syscall.Chdir("/")

	// repeating previous commands from run function now that we have a proper namespace set up
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// mounting the "image"
	syscall.Mount("proc", "proc", "proc", 0, "")
	// other files can be mounted here using the same command as above

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	// cleanup
	syscall.Unmount("proc", 0)
}

func controlGroups() {
	cgroups := LoadEnv("CGROUPS")
	pids := filepath.Join(cgroups, "pids")

	os.Mkdir(filepath.Join(pids, "container"), 0755)

	ioutil.WriteFile(filepath.Join(pids, "container/pids.max"), []byte("20"), 0700)

	// removing created control group once the container is closed
	ioutil.WriteFile(filepath.Join(pids, "container/notify_on_release"), []byte("1"), 0700)
	ioutil.WriteFile(filepath.Join(pids, "container/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
}
