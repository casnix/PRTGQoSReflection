# PRTGQoSReflection
This is a reproduction of github.com/PaesslerAG/QoSReflect written in Golang with my style of build scripts.  The purpose of this is to have a standalone executable
that can be built on any OS and provides the same functionality of PaesslerAG/QoSReflect without depending on Python.

## Execution prerequisites
None, really.  

## Build prerequisites
* Golang 1.18.1
* Bash environment (for build scripts) -- can be on Linux or WSL, or something like WSL.

## Installation
* [OPTIONAL]create a file called "qosreflect.conf" with the following contents:

host:All
port:50000
replyip:None

The script can now be called with parameters to allow several instances running. Just type qosreflect<.exe> --help to see all parameters. Example call below:

qosreflect.exe --port 50000 --host All

Additional parameters are optional. You can still use a config file, then please use parameter --conf to provide the path.

When "host" is set to "All" the script will try to bind to every available interface. Change to IP of an interface to make the script bind to a special interface.
Set "port" to the same one set up in PRTG.
If an IP is specified in "replyip" the script will only process UDP packets from this IP.

## Debugging
To debug whats going on call the script with the additional parameter -d or --debug