package main

import (
	"log"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

func BashHandler(serverName, background string) ssh.Handler {
	return func(s ssh.Session) {
		ctx := s.Context()
		simulator, err := NewSimulator(BashMode, serverName, background)
		if err != nil {
			log.Fatal(err)
		}

		t := term.NewTerminal(s, "> ")

		for {
			line, err := t.ReadLine()
			if err != nil {
				log.Println("err:", err)
				return
			}
			if line == "" {
				continue
			}
			result, err := simulator.Simulate(ctx, line)
			if err != nil {
				log.Fatal(err)
			}
			s.Write([]byte(result + "\n"))
		}
	}
}
