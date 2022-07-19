# Rootless Containers
A rootless container system written in Go following along Liz Rice's presentation at the goto conference.

### Setup: 
* Create a filesystem to use a container image, and specify the path to the filesystem in `.env` under `CONTAINER_ROOT`
* Create a control group system, or use the one already on your machine. If you'd like to use the one already on your computer, set `CGROUPS` in `.env` to `/sys/fs/cgroup`

### Running a container:
* To run the container, run `go run main.go run <COMMAND>` (meant to replicate docker run IMAGE COMMAND)
* To enter the shell of the container, run `go run main.go run /bin/bash`
