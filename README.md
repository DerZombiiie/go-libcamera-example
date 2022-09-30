# This is a small example showing how you could and arguably shouln't read data from stdout of some process.

...or how to get some kind of image from libcamera-still in go.

This is a proof of concept!

Usage: 

- `go run libcamera.go`

- note the PID (first line) should be sth like 17731

- `kill -SIGUSR1 <PID>` - signals libcamera-still to take a picture

- wait.. after a second or so the programm should tell you a filename where the image is saved
  name is out-<unixtime>.jpg

this example is under the MIT license, see LICENSE
