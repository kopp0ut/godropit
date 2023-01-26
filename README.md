# godropit

Simple Dropper Automation tool just enough to automate using some common go payloads from the excellent examples of [Ne0nd0g](https://github.com/Ne0nd0g/go-shellcode)

## Install
To use gohide it you will need to install go, as we'll need the go compiler to compile our droppers. You'll also want mingw-w64 for any payload using C, e.g. DLL droppers.

For Linux
```sh
sudo apt update && sudo apt install golang mingw-w64
git clone https://gitlab.alphaweasel.com/rt/purple-teaming/godropit
```


## Usage
Once downloaded it's fairly simple, you can either compile the binary and run it that way, or use it with "go run". Right now it only accepts baked in templates, but it will output the src to "-o outputdir". This is where all the go source files will be placed, including encrypted shellcode, env variables used and go.mod/go.sum etc. 

```sh
.\godropit new <local|child|remote> -i <exe/binfile> -n <outputfilename> -o <output directory>
```
[![asciicast](https://asciinema.org/a/C6y3GDaCQEBZn6NXSWg6QrlPx.svg)](https://asciinema.org/a/C6y3GDaCQEBZn6NXSWg6QrlPx)

## Supported Functionality
* Domain keying
* Local, remote and child process execution methods.
* Automatic AES shellcode encryption
* 

## Tips

### Test Shellcode
**-i flag**
If you are unsure about 

### Execution delays
**-t flag**
This takes the time in seconds, and adds that as a delay prior to execution and decryption. Can be useful for distancing the exec of the binary from injection activity by a margin.
If not specified, this will only delay by 1 second. 

### Domain Keying
**-d flag**
Domain keying is supported, it will check box the go hostname (which often includes domain\\) and the USERDNSDOMAIN env var.
This uses strings.Contain() and is not case sensitive, so you can put as much or as little of the domname in to this.

There are plans to support more, but I suggest you write your own bespoke payloads for most scenarios - this is just for quick wins against simple AV.

## Planned Features:

* Support for using [garble](https://github.com/burrowers/garble) against your binary when you are done.
* Support for custom exports in dlls.
* Better string and import obfuscation, e.g. using []byte or encryption rather than plain strings.
* LocalKCT execution (some implemented) - https://github.com/aahmad097/AlternativeShellcodeExec
* PPID Spoofing
* ProcByName
* Custom Export names
