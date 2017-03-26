package exec

import (
	"context"
	"errors"
	"io"
	"os/user"
)

type Command interface {
	Start(ctx context.Context) error
	Stop() error
}

type Interpreter interface {
	Cmd(config Config, ctx context.Context) (Command, error)
}

var interpreters map[string]Interpreter

// Run config
type Config struct {
	// Command and arguments to run
	Cmd []string
	// Extra arguments
	Args []string
	// Inline script
	Script string
	// Environment
	Env Environ
	// Working dir
	WorkDir string
	// Interpreter and arguments
	Interpreter []string
	User        *user.User
	Stdout      io.Writer
	Stderr      io.Writer
	Stdin       io.Reader
}

func Register(name string, i Interpreter) {
	interpreters[name] = i
}

func init() {
	interpreters = make(map[string]Interpreter)
	interpreters["shell"] = &shell{}
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

	intp := interpreters["shell"]
	if conf.Interpreter == nil {
		intp = interpreters["shell"]
	} else if i, o := interpreters[conf.Interpreter[0]]; o {
		intp = i
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
