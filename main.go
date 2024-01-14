package main

import (
	"io"
	"os"
)

func main() {
	err := launch(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func launch(in io.ReadCloser, out io.WriteCloser) error {
	g, err := newGovimPlugin(in, out)
	if err != nil {
		return err
	}
	if err := g.init(); err != nil {
		return err
	}
	return g.loop()
}
