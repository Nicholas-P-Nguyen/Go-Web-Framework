# Setting up the Development Environment

In this section, we you will get your development environment set up
with the necessary language compilers, libraries, and tools to complete
assignment 1 and future assignments & precepts.

## Development Environment for Assignments and Precepts

# Setting up the Development Environment

In this section, we you will get your development environment set up
with the necessary language compilers, libraries, and tools to complete
assignment 1 and future assignments & precepts.

## Development Environment for Assignments and Precepts

The assignments and precepts in this course assume that you have a few tools and
libraries available in your development environment. Any way you get them is
fine, and they should generally be relatively easy to find and install on most
Linux distributions, BSDs, macOS, and Windows using Windows Subsystem for
Linux (WSL).  

Note: before installing, you may go down to the "Verifying you have the right tools installed" section first to see what tools
you already have; chances are you've already installed some of them (like git) locally.

  * [Git](https://git-scm.org)
  * [Go](https://go.dev) 
  * [SQLite](https://www.sqlite.org/index.html) version 3
  * [curl](https://curl.se)
  * [gcc/g++](https://gcc.gnu.org/) OR [clang/clang++](https://clang.llvm.org/) compilers
    for C & C++ (other C/C++ compiler may or may not work).
  * [make](https://www.gnu.org/software/make/)

In addition you should have either Firefox or Chrom[e|ium] installed for one of the precepts where we will inspect HTTP traffic using those browsers. Recent versions of the Edge browser which are based on Chromium and
include Chromium's developer tools will likely also work.

Note: All commands below will start with a '$' character; be sure to exclude it from the actual command you're running.


## Installing on Linux

Most recent Linux distributions should have the necessary tools available from
their respective package managers. Some Long-Term-Support (LTS) versions may
have an older version of Go, in which case you can use the instructions on the
Go website to install a newer version of Go.

Below are instructions for a couple popular distributions. You will typically
need root permissions (i.e. prefix the commands with `sudo` or login as root) to
install packages:

### Debian (Backports/Testing/Unstable) & Ubuntu (Rolling)

``` sh
$ apt-get install git curl build-essential golang sqlite3
```

### Arch Linux

``` sh
$ pacman -S git go sqlite3 curl base-devel
```

### Alpine Linux

``` sh
$ apk add git go build-base sqlite curl
```

## Installing on macOS

First, you must have the [XCode Command Line Tools](https://developer.apple.com/downloads/index.action) from Apple installed. These already include a C/C++ compiler and `make`.

You can obtain the remaining tools from a number of package managers available for macOS (e.g. Homebrew, MacPorts), or you can also download and install universal packages from each tool's respective website.

A simple way using [MacPorts](https://www.macports.org/) is with the following command:

``` sh
$ sudo port install git go curl sqlite3
```

Or with [homebrew](https://brew.sh/):
``` sh
$ brew install git go curl sqlite3
```

## Installing on Windows under WSL

To work in a Windows environment, we recommend using Windows Subsystem for
Linux, which provides a Linux environment alongside Windows. By default, WSL
provides a Ubuntu distribution of Linux. Follow the installation instructions
for WSL from the [official
documentation](https://docs.microsoft.com/en-us/windows/wsl/install), then use
the instructions above for Linux to install the tools.

You _may_ be able to complete assignments with these tools installed on
    Windows _without_ WSL, but your mileage may vary and course staff may not be
    able to offer help with such a setup. We strongly recommend enabling and
    using WSL if you are on Windows.

Most likely, you will be using the default Ubuntu flavor of WSL, in which case
the following command will install all the necessary tools:

``` sh
$ sudo apt-get update
$ sudo apt-get install git curl build-essential golang sqlite3
```


## Verifying you have the right tools installed

Once you've installed the tools, running the
following commands in a shell should result in similar output (slight variations
are expected on different operating systems). Also, the exact output and versions
of the tools may be different, so don't worry about a version mismatch between what
you have locally; as long as they are (relatively) up to date then you should be fine.

### Git

``` sh
$ git version
git version 2.34.1
```

### Go

``` sh
$ go version
go version go1.18.5 linux/amd64
```

### SQLite 3

``` sh
$ echo "select (316 + 1021)" | sqlite3
1337
```

### curl

``` sh
$ curl -I https://www.google.com
HTTP/2 200
server: nginx
date: Wed, 26 Jan 2022 16:41:20 GMT
content-type: text/html
content-length: 3859
etag: "xsmg9fnyjsw7c3qidhrrb62g071vfnls"
accept-ranges: bytes
...
```

### C++ compiler

``` sh
$ g++ --version
g++ (GCC) 10.3.0
Copyright (C) 2020 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
```

OR

``` sh
$ clang --version
clang version 7.1.0 (tags/RELEASE_710/final)
Target: x86_64-unknown-linux-gnu
Thread model: posix
InstalledDir: /nix/store/ass1sf1bx07qvlrg02nymxnyzp1cpxz7-clang-7.1.0/bin
```

### Make

``` sh
$ make --version
GNU Make 4.3
Built for x86_64-pc-linux-gnu
Copyright (C) 1988-2020 Free Software Foundation, Inc.
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
```
