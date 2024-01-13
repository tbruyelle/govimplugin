package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

func main() {
	err := launch(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func launch(in io.ReadCloser, out io.WriteCloser) error {
	g := govimPlugin{
		in:  json.NewDecoder(in),
		out: json.NewEncoder(out),
		log: slog.Default(),
	}

	for {
		g.log.Debug("run: waiting to read JSON Message")
		id, msg := g.readJSONMsg()
		g.log.Info("readJSONMsg:", "id", id, "msg", msg)
		args := g.parseJSONArgSlice(msg)
		typ := g.parseString(args[0])
		args = args[1:]
		g.log.Info("parsed:", "type", typ, "args", args)
	}
}
