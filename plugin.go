package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

type govimPlugin struct {
	in  *json.Decoder
	out *json.Encoder

	// outLock synchronises access to out to ensure we have non-overlapping
	// sending of messages
	outLock sync.Mutex
	log     *slog.Logger
}

// readJSONMsg is a low-level protocol primitive for reading a JSON msg sent by Vim.
// There is more structure to the messages that we receive, hence we can be
// more specific in our return type. See
// https://vimhelp.org/channel.txt.html#channel-use for more details.
func (g *govimPlugin) readJSONMsg() (int, json.RawMessage) {
	var msg [2]json.RawMessage
	if err := g.in.Decode(&msg); err != nil {
		if err == io.EOF {
			// explicitly setting underlying here
			panic(errProto{underlying: err})
		}
		g.errProto("failed to read JSON msg: %v", err)
	}
	i := g.parseInt(msg[0])
	return i, msg[1]
}

// parseJSONArgSlice is a low-level protocol primitive for parsing a slice of
// raw encoded JSON values
func (g *govimPlugin) parseJSONArgSlice(m json.RawMessage) []json.RawMessage {
	var i []json.RawMessage
	g.decodeJSON(m, &i)
	return i
}

// parseString is a low-level protocol primtive for parsing a string from a
// raw encoded JSON value
func (g *govimPlugin) parseString(m json.RawMessage) string {
	var s string
	g.decodeJSON(m, &s)
	return s
}

// parseInt is a low-level protocol primtive for parsing an int from a
// raw encoded JSON value
func (g *govimPlugin) parseInt(m json.RawMessage) int {
	var i int
	g.decodeJSON(m, &i)
	return i
}

// decodeJSON is a low-level protocol primitive for decoding a JSON value.
func (g *govimPlugin) decodeJSON(m json.RawMessage, i interface{}) {
	err := json.Unmarshal(m, i)
	if err != nil {
		g.errProto("failed to decode JSON into type %T: %v", i, err)
	}
}

func (g *govimPlugin) errProto(format string, args ...interface{}) {
	panic(errProto{
		underlying: fmt.Errorf(format, args...),
	})
}

type errProto struct {
	underlying error
}

func (e errProto) Error() string {
	return fmt.Sprintf("protocol error: %v", e.underlying)
}

// sendJSONMsg is a low-level protocol primitive for sending a JSON msg that will be
// understood by Vim. See https://vimhelp.org/channel.txt.html#channel-use
func (g *govimPlugin) sendJSONMsg(p1, p2 any, ps ...any) {
	msg := []any{p1, p2}
	msg = append(msg, ps...)
	// TODO: could use a multi-writer here
	logMsg, err := json.Marshal(msg)
	if err != nil {
		g.errProto("failed to create log message: %v", err)
	}
	g.log.Info("sendJSONMsg: %s\n", logMsg)
	g.outLock.Lock()
	defer g.outLock.Unlock()
	if err := g.out.Encode(msg); err != nil {
		panic(err)
	}
}
