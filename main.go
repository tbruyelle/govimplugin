package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	err := launch(os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func launch(in io.ReadCloser, out io.WriteCloser) error {
	f, err := os.Create(filepath.Join(os.TempDir(), "govimplugin.log"))
	if err != nil {
		return err
	}
	defer f.Close()

	g := govimPlugin{
		in:  json.NewDecoder(in),
		out: json.NewEncoder(out),
		log: slog.New(slog.NewJSONHandler(f, nil)),
	}

	for {
		g.log.Info("run: waiting to read JSON Message")
		id, msg := g.readJSONMsg()
		g.log.Info("readJSONMsg:", "id", id, "msg", msg)
		args := g.parseJSONArgSlice(msg)
		typ := g.parseString(args[0])
		args = args[1:]
		g.log.Info("parsed:", "type", typ, "args", args)
	}
}
