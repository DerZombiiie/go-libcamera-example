package main

import (
	"github.com/DerZombiiie/go-libcamera-example/collector"

	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("libcamera-still",
		"-o-",   // stdout
		"-s",    // listen on signals like SIGUSR1
		"-t0",   // no max runtime
		"-ejpg", // jpeg encoding
	)

	// create writecollector with 250ms timeout
	cmd.Stdout = collector.MakeCollector(time.Millisecond*250, func(buf *bytes.Buffer) {
		file := fmt.Sprintf("out-%d.jpg", time.Now().Unix())

		f, err := os.Create(file)
		if err != nil {
			log.Printf("Error creating file '%s': %s \n", file, err)
			return
		}

		defer f.Close()

		io.Copy(f, buf)
		log.Printf("Created file '%s' \n", file)
	})

	cmd.Stderr = os.Stderr // passthrough error

	// start process
	if err := cmd.Start(); err != nil {
		panic(err)
	}

	// log pid
	log.Printf("PID: %d \n", cmd.Process.Pid)

	// wait for process to get killed
	cmd.Wait()
}
