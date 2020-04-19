# dls
docker ls - list files within a container

This program is used to view what files and directories are inside of a running docker container.  In your `Dockerfile`,
you will need to `ADD` or `COPY` a **statically linked** `dls` binary.

## Usage

```
dls: Get file info for a directory
usage: dls [directory]
  -a	show all files, including .git, dev, proc, and sys
  -b	show in bare format, eg: no tables
  -e	show file/directory errors
  -t	show total file size of all files
  -v	show version and then exit

Current directory is the default if no other directory is given on cmd-line
```

## Example

* dls.exe -t

```
+-----------+---------------------+------+------------------------+
|   SIZE    |      MOD TIME       | TYPE | NAME (FILES:11 DIRS:2) |
+-----------+---------------------+------+------------------------+
|      1086 | 2020-04-16 15:21:34 | F    | LICENSE                |
|        50 | 2020-04-16 15:21:34 | F    | README.md              |
|      4525 | 2020-04-16 23:12:31 | F    | cmd.go                 |
|   3209728 | 2020-04-16 23:12:39 | F    | dls.exe                |
|        81 | 2020-04-16 23:09:12 | F    | osver_darwin.go        |
|        82 | 2020-04-16 23:10:13 | F    | osver_freebsd.go       |
|        80 | 2020-04-16 23:09:47 | F    | osver_linux.go         |
|       964 | 2020-04-16 23:05:04 | F    | osver_windows.go       |
|         0 | 2020-04-16 23:15:46 | D    | test1\                 |
|       351 | 2020-04-16 23:15:41 | F    | test1\file1.txt        |
|       856 | 2020-04-16 23:15:46 | F    | test1\file2.txt        |
|         0 | 2020-04-16 23:16:03 | D    | test2\                 |
|       581 | 2020-04-16 23:16:04 | F    | test2\file3.txt        |
| --------- | ------------------- | -    | ---------------        |
|           |             3218384 |      | (total size)           |
|           |                3.07 |      | (MB total size)        |
+-----------+---------------------+------+------------------------+
```

## Docker

### Static Binary Complication

| Platform | Command
----------|-----
| windows | go build -tags netgo -ldflags "-extldflags -static"
| linux/bsd | go build -tags netgo -ldflags '-extldflags "-static" -s -w'
| macos | go build -ldflags '-s -extldflags "-sectcreate __TEXT __info_plist Info.plist"'
| android | go build -ldflags -s

**NOTE:** *I have not been able to test all of these*

### Windows
* You will need to use `nanoserver` with a matching OS version: [Docker Image List](https://hub.docker.com/_/microsoft-windows-nanoserver)
* * Running `dls -v` will return your specific `OS build`
* Example for Windows 10 LTSC Build 1809:
* * `docker pull mcr.microsoft.com/windows/nanoserver:10.0.17763.1158`

### Linux
* Remember to create a static version of the binary
* `dls` can be started from the [Docker Scratch Image](https://hub.docker.com/_/scratch)

### References
* [Building Minimal Docker Containers for Go Applications](https://rollout.io/blog/building-minimal-docker-containers-for-go-applications/)
* ["Distroless" Docker Images](https://github.com/GoogleContainerTools/distroless)
* [Statically compiling Go programs](https://www.arp242.net/static-go.html)
