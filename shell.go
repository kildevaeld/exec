package exec

import (
	"context"
	"errors"
	"os/exec"
	"os/user"
	"strconv"

	"syscall"
)

type WrapCommand struct {
	cmd     *exec.Cmd
	running bool
}

func (self *WrapCommand) Start(ctx context.Context) error {
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

type shell struct{}

func (self *shell) Cmd(config Config, ctx context.Context) (Command, error) {

	var cmd *exec.Cmd

	var path string
	var args []string

	if len(config.Interpreter) > 0 {
		path = config.Interpreter[0]
		args = config.Interpreter[1:]
	}

	if config.Script != "" {
		if path == "" {
			path = "sh"
			if len(args) == 0 {
				args = []string{"-c"}
			}
		}
		args = append(args, config.Script, "--")
	} else if len(config.Cmd) > 0 {
		if path == "" {
			path = config.Cmd[0]
			args = append(args, config.Cmd[1:]...)
		} else {
			args = append(args, config.Cmd...)
		}

	} else {
		return nil, errors.New("No command or no script")
	}

	if config.Args != nil {
		args = append(args, config.Args...)
	}

	cmd = exec.CommandContext(ctx, path, args...)

	cmd.Stderr = config.Stderr
	cmd.Stdout = config.Stdout
	cmd.Stdin = config.Stdin

	cmd.Dir = config.WorkDir
	cmd.Env = config.Env

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if config.User != nil {
		var (
			err   error
			creds *syscall.Credential
		)

		if creds, err = getCredentials(config.User); err != nil {
			return nil, err
		}

		cmd.SysProcAttr.Credential = creds

	}

	return &WrapCommand{cmd, false}, nil
}

func getCredentials(user *user.User) (*syscall.Credential, error) {
	var (
		i int
		g int
		e error
	)

	if i, e = strconv.Atoi(user.Uid); e != nil {
		return nil, e
	}

	if g, e = strconv.Atoi(user.Gid); e != nil {
		return nil, e
	}

	groupids, err := user.GroupIds()
	if err != nil {
		return nil, err
	}

	var gids []uint32
	for _, gid := range groupids {
		if g, e = strconv.Atoi(gid); e != nil {
			return nil, e
		} else {
			gids = append(gids, uint32(g))
		}
	}

	return &syscall.Credential{
		Uid:    uint32(i),
		Gid:    uint32(g),
		Groups: gids,
	}, nil
}
