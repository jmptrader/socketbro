package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// ===== GLOBALS ===================================================================================

var listeners []net.Listener

// ===== FUNCTIONS =================================================================================

func newUnixSocketListener(path string) net.Listener {
	fi, err := os.Stat(path)
	if nil == err {
		if fi.Mode() & os.ModeSocket != os.ModeSocket {
			log.Fatalln("Refusing to remove existing non-socket file:", path)
		}

		err := syscall.Unlink(path)
		if err != nil {
			log.Fatalf("Failed to remove old socket file %s: %s\n", path, err)
		}
	}

	listener, err := net.Listen("unix", path)
	if err != nil {
		log.Fatalln("net.Listen:", err)
	}

	log.Println("Socket created at", path)

	return listener
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for _ = range c {
		log.Println("Received signal, cleaning up")
		for _, listener := range listeners {
			err := listener.Close()
			if err != nil {
				log.Fatalln("net.Listener.Close:", err)
			}
		}

		break
	}

	os.Exit(0)
}

func main() {
	var paths string

	flag.StringVar(&paths, "paths", "", "CSV of unix sockets, e.g., /tmp/1.sock,/tmp/2.sock")
	flag.Parse()

	if 0 == len(paths) {
		log.Fatalln("You must specify --paths")
	}

	for _, path := range strings.Split(paths, ",") {
		listeners = append(listeners, newUnixSocketListener(path))
	}

	go handleSignals()

	for _ = range time.Tick(1 * time.Second) {
		// do nothing, forever
	}
}
