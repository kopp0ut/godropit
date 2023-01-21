# godropit

Simple Dropper Automation tool just enough to automate using some common go payloads from the excellent examples of [Ne0nd0g](https://github.com/Ne0nd0g/go-shellcode)

## Install
To use gohide it you will need to install go, as we'll need the go compiler to compile our droppers. You'll also want mingw-w64 for any payload using C, e.g. DLL droppers.

For linux:
```sh
sudo apt update && sudo apt install golang
git clone https://gitlab.alphaweasel.com/rt/purple-teaming/godropit
```

## Usage
Once downloaded it's fairly simple, you can either compile the binary and run it that way, or use it with "go run". However, please note that the tool needs to be run in a folder with the cmd/ payload directory as this is where it looks for the available payloads. Any other go payloads you put in there will show as options to use, but please note they must match the template (BasicObf.go) to actually work.

```sh
.\godropit -i <exe/binfile> -aeskey <AES Key String> -o <output executable>
```
[![asciicast](https://asciinema.org/a/nyFglEvadib7FNVGwTzK62cb4.png)](https://asciinema.org/a/nyFglEvadib7FNVGwTzK62cb4)

## Notes
Currently this will only generate 64 bit exe's, which is all I plan to use this for. If you need 32bit, you can tweak the code (generate.go), but I don't support 32-bit myself so please don't come to me if your payload of choice does not work.

There are plans to support more, but I suggest you write your own bespoke payloads for most scenarios - this is just for quick wins against simple AV.

## Planned Features:

* Support for using [garble](https://github.com/burrowers/garble) against your binary when you are done.
* Support for custom exports in dlls.
