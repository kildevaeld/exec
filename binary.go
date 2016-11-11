package exec

import (
	"errors"
	"os/exec"
)

type WrapCommand struct {
	cmd     *exec.Cmd
	running bool
}

func (self *WrapCommand) Start() error {
	if self.running {
		return errors.New("already running")
	}
	self.running = true
	return self.cmd.Run()
}

func (self *WrapCommand) Stop() error {
	if !self.running {
		return nil
	}

	if self.cmd.Process != nil {
		return self.cmd.Process.Kill()
	}
	return nil
}

type binary struct{}

func (self *binary) Cmd(config Config) Command {

	var cmd *exec.Cmd

	if len(config.Cmd) > 1 {
		cmd = exec.Command(config.Cmd[0], config.Cmd[1:]...)
	} else {
		cmd = exec.Command(config.Cmd[0])
	}

	cmd.Stderr = config.Stderr
	cmd.Stdout = config.Stdout

	cmd.Dir = config.WorkDir
	cmd.Env = config.Env

	return &WrapCommand{cmd, false}
}
