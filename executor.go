package exec

import (
	"context"
	"errors"
	"io"
)

type Command interface {
	Start(ctx context.Context) error
	Stop() error
}

type Interpreter interface {
	Cmd(config Config, ctx context.Context) (Command, error)
}

var interpreters map[string]Interpreter

type Config struct {
	Cmd         []string
	Args        []string
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

func (self *Executor) Start(ctx context.Context) (err error) {
	if self.cmd != nil {
		return errors.New("already started")
	}
	conf := self.config

	var intp Interpreter
	if conf.Interpreter == nil || len(conf.Interpreter) == 0 {
		intp = &shell{}
	} else if len(conf.Interpreter) > 0 {
		if i, o := interpreters[conf.Interpreter[0]]; o {
			intp = i
		} else {
			intp = &shell{}
		}
	}

	if self.cmd, err = intp.Cmd(conf, ctx); err != nil {
		return err
	}

	return self.cmd.Start(ctx)
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
