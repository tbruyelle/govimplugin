package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type govimPlugin struct {
	in  *json.Decoder
	out *json.Encoder

	// outLock synchronises access to out to ensure we have non-overlapping
	// sending of messages
	outLock sync.Mutex
	log     *slog.Logger

	callbackRespsLock sync.Mutex
	callVimNextID     int
	callbackResps     map[int]chan callbackResp

	callbacks map[string]callback
}

type callback func(name string, args ...string) error

// callbackResp is the container for a response from a call to callVim. If the
// call does not result in a value, e.g. ChannelEx, then val will be nil
type callbackResp struct {
	errString string
	val       json.RawMessage
}

func newGovimPlugin(in io.Reader, out io.Writer) (*govimPlugin, error) {
	f, err := os.Create(filepath.Join(os.TempDir(), "govimplugin.log"))
	if err != nil {
		return nil, err
	}
	// defer f.Close()

	g := &govimPlugin{
		in:            json.NewDecoder(in),
		out:           json.NewEncoder(out),
		log:           slog.New(slog.NewJSONHandler(f, nil)),
		callbacks:     make(map[string]callback),
		callbackResps: make(map[int]chan callbackResp),
	}
	go func() {
		err := g.loop()
		if err != nil {
			panic(err)
		}
	}()
	return g, nil
}

func (g *govimPlugin) init() error {
	g.mustDefineCommand("TestCmd", g.testCmdCb)
	return nil
}

func (g *govimPlugin) testCmdCb(name string, args ...string) error {
	g.log.Info("testCmdCb", "name", name, "args", args)
	return nil
}

func (g *govimPlugin) loop() error {
	for {
		g.log.Info("run: waiting to read JSON Message")
		id, msg := g.readJSONMsg()
		g.log.Info("readJSONMsg:", "id", id, "msg", msg)
		args := g.parseJSONArgSlice(msg)
		typ := g.parseString(args[0])
		args = args[1:]
		g.log.Info("parsed:", "type", typ, "args", args)
		switch typ {
		case "callback":
			// This case is a "return" from a call to callVim. Format of args
			// will be [id, [string, val]]
			id := g.parseInt(args[0])
			resp := g.parseJSONArgSlice(args[1])
			msg := g.parseString(resp[0])
			var val json.RawMessage
			if len(resp) == 2 {
				val = resp[1]
			}
			toSend := callbackResp{
				errString: msg,
				val:       val,
			}
			g.callbackRespsLock.Lock()
			ch, ok := g.callbackResps[id]
			delete(g.callbackResps, id)
			g.callbackRespsLock.Unlock()
			if !ok {
				g.errProto("run: received response for callback %v, but not response chan defined", id)
			}
			select {
			case ch <- toSend:
			default:
				g.log.Error("nobody's listening to chan")
			}
			/*
				switch ch := ch.(type) {
				case scheduledCallback:
					g.eventQueue.Add(func() error {
						select {
						case ch <- toSend:
						case <-g.tomb.Dying():
							return ErrShuttingDown
						}
						return nil
					})
				case unscheduledCallback:
					g.tomb.Go(func() error {
						select {
						case ch <- toSend:
						case <-g.tomb.Dying():
							return tomb.ErrDying
						}
						return nil
					})
				default:
					panic(fmt.Errorf("unknown type of callback responser: %T", ch))
				}
			*/
		}
	}
}

func (g *govimPlugin) mustDefineCommand(name string, cb callback, attrs ...CommAttr) {
	err := g.defineCommand(name, cb, attrs...)
	if err != nil {
		panic(err)
	}
}

func (g *govimPlugin) defineCommand(name string, cb callback, attrs ...CommAttr) error {
	if _, ok := g.callbacks[name]; ok {
		return fmt.Errorf("command %q already defined", name)
	}
	g.log.Info("define command", "name", name)
	var nargsFlag *NArgs
	var rangeFlag *Range
	var rangeNFlag *RangeN
	var countNFlag *CountN
	var completeFlag *CommAttr
	genAttrs := make(map[CommAttr]bool)
	for _, iattr := range attrs {
		iattr := iattr
		switch attr := iattr.(type) {
		case NArgs:
			switch attr {
			case NArgs0, NArgs1, NArgsZeroOrMore, NArgsZeroOrOne, NArgsOneOrMore:
			default:
				return fmt.Errorf("unknown NArgs value")
			}
			if nargsFlag != nil && attr != *nargsFlag {
				return fmt.Errorf("multiple nargs flags")
			}
			nargsFlag = &attr
		case Range:
			switch attr {
			case RangeLine, RangeFile:
			default:
				return fmt.Errorf("unknown Range value")
			}
			if rangeFlag != nil && *rangeFlag != attr || rangeNFlag != nil {
				return fmt.Errorf("multiple range flags")
			}
			if countNFlag != nil {
				return fmt.Errorf("range and count flags are mutually exclusive")
			}
			rangeFlag = &attr
		case RangeN:
			if rangeNFlag != nil && *rangeNFlag != attr || rangeFlag != nil {
				return fmt.Errorf("multiple range flags")
			}
			if countNFlag != nil {
				return fmt.Errorf("range and count flags are mutually exclusive")
			}
			rangeNFlag = &attr
		case CountN:
			if countNFlag != nil && *countNFlag != attr {
				return fmt.Errorf("multiple count flags")
			}
			if rangeFlag != nil || rangeNFlag != nil {
				return fmt.Errorf("range and count flags are mutually exclusive")
			}
			countNFlag = &attr
		case Complete:
			if completeFlag != nil && *completeFlag != attr {
				return fmt.Errorf("multiple complete flags")
			}
			completeFlag = &iattr
		case CompleteCustom:
			if completeFlag != nil && *completeFlag != attr {
				return fmt.Errorf("multiple complete flags")
			}
			completeFlag = &iattr
		case CompleteCustomList:
			if completeFlag != nil && *completeFlag != attr {
				return fmt.Errorf("multiple complete flags")
			}
			completeFlag = &iattr
		case GenAttr:
			switch attr {
			case AttrBang, AttrRegister, AttrBuffer, AttrBar:
				genAttrs[attr] = true
			default:
				return fmt.Errorf("unknown GenAttr value")
			}
		}
	}
	attrMap := make(map[string]interface{})
	if nargsFlag != nil {
		attrMap["nargs"] = nargsFlag.String()
	}
	if rangeFlag != nil {
		attrMap["range"] = rangeFlag.String()
	}
	if rangeNFlag != nil {
		attrMap["range"] = rangeNFlag.String()
	}
	if countNFlag != nil {
		attrMap["count"] = countNFlag.String()
	}
	if completeFlag != nil {
		attrMap["complete"] = (*completeFlag).String()
	}
	if len(genAttrs) > 0 {
		var attrs []string
		for k := range genAttrs {
			attrs = append(attrs, k.String())
		}
		sort.Strings(attrs)
		attrMap["general"] = attrs
	}
	args := []interface{}{name, attrMap}
	ch := make(chan callbackResp)
	// err = g.DoProto(func() error {
	err := g.callVim(ch, "command", args...)
	// })
	return g.handleChannelError(ch, err, "failed to define %q in Vim: %v", name)
}

func (g *govimPlugin) handleChannelError(ch chan callbackResp, err error, format string, args ...interface{}) error {
	_, err = g.handleChannelValueAndError(ch, err, format, args...)
	return err
}

func (g *govimPlugin) handleChannelValueAndError(ch chan callbackResp, err error, format string, args ...interface{}) (json.RawMessage, error) {
	if err != nil {
		return nil, err
	}
	args = append([]interface{}{}, args...)
	resp := <-ch
	if resp.errString != "" {
		args = append(args, resp.errString)
		return nil, fmt.Errorf(format, args...)
	}
	return resp.val, nil
	/*
		select {
		case <-g.tomb.Dying():
			panic(ErrShuttingDown)
		case resp := <-ch:
			if resp.errString != "" {
				args = append(args, resp.errString)
				return fmt.Errorf(format, args...)
			}
		}
	*/
}

// callVim is a low-level protocol primitive for making a call to the
// channel defined handler in Vim. The Vim handler switches on typ. The Vim
// handler does not return a value, instead it will acknowledge success by
// sending a zero-length string.
func (g *govimPlugin) callVim(ch chan callbackResp, typ string, vs ...interface{}) error {
	g.callbackRespsLock.Lock()
	id := g.callVimNextID
	g.callVimNextID++
	g.callbackResps[id] = ch
	g.callbackRespsLock.Unlock()
	args := []interface{}{id, typ}
	args = append(args, vs...)
	g.sendJSONMsg(0, args)
	return nil
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
	g.log.Info("sendJSONMsg:", "msg", logMsg)
	g.outLock.Lock()
	defer g.outLock.Unlock()
	if err := g.out.Encode(msg); err != nil {
		panic(err)
	}
}
