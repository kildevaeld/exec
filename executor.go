package exec

import (
	"errors"
	"io"
)

type Environ []string

type Command interface {
	Start() error
	Stop() error
}

type Interpreter interface {
	Cmd(config Config) Command
}

var interpreters map[string]Interpreter

type Config struct {
	Cmd         []string
	Script      string
	Env         Environ
	WorkDir     string
	Interpreter []string
	Stdout      io.Writer
	Stderr      io.Writer
	Stdin       io.Reader
}

func Register(name string, i Interpreter) {
	interpreters[name] = i
}

func init() {
	interpreters = make(map[string]Interpreter)
}

type Executor struct {
	config Config
	cmd    Command
}

func (self *Executor) Start() error {
	if self.cmd != nil {
		return errors.New("already started")
	}
	conf := self.config

	var intp Interpreter
	if conf.Interpreter == nil || len(conf.Interpreter) == 0 {
		intp = &binary{}
	} else if len(conf.Interpreter) > 0 {
		if i, o := interpreters[conf.Interpreter[0]]; o {
			intp = i
		} else {
			intp = &binary{}
		}
	}

	self.cmd = intp.Cmd(conf)

	return self.cmd.Start()
}

func (self *Executor) Stop() error {
	if self.cmd == nil {
		return errors.New("Not started")
	}
	return self.cmd.Stop()
}

func New(config Config) *Executor {
	return &Executor{config, nil}
}
