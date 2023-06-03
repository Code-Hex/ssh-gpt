package main

import (
	"flag"
	"log"

	"github.com/gliderlabs/ssh"
)

const (
	netflixServer     = "ssh.netflix.net"
	netflixBackground = "This is exposed by Netflix to share exclusive materials with their fans, includings a bunch of folders, movies, txt files and more."
)

func main() {
	serverName := flag.String("s", netflixServer, "hostname of fictitious ssh server")
	background := flag.String("bg", netflixBackground, "background for fictitious ssh server")
	flag.Parse()
	srv := &ssh.Server{
		Addr:    ":2222",
		Handler: BashHandler(*serverName, *background),
		SubsystemHandlers: map[string]ssh.SubsystemHandler{
			"sftp": ssh.SubsystemHandler(SFTPHandler(*serverName, *background)),
		},
	}

	log.Printf("starting ssh server as %s on port 2222...", *serverName)
	log.Fatal(srv.ListenAndServe())
}
