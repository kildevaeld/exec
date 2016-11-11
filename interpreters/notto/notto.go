package notto

import (
	"context"
	"path/filepath"

	"github.com/kildevaeld/exec"
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules"
)

type NottoCommand struct {
	vm     *notto.Notto
	isFile bool
	cmd    string
}

func (self *NottoCommand) Start(ctx context.Context) error {

	var err error
	if self.isFile {
		_, err = self.vm.Run(self.cmd, filepath.Dir(self.cmd))
	} else {
		_, err = self.vm.RunScript(self.cmd, self.vm.ProcessAttr().Cwd)
	}

	return err
}

func (self *NottoCommand) Stop() error {

	return nil
}

type NodeInterpreter struct {
}

func (self *NodeInterpreter) Cmd(config exec.Config, ctx context.Context) (exec.Command, error) {

	vm := notto.New()

	if err := modules.Define(vm); err != nil {
		return nil, err
	}
	var args []string

	n := &NottoCommand{vm, false, config.Script}

	if len(config.Cmd) > 0 {
		args = config.Cmd[1:]
		n.isFile = true
		n.cmd = config.Cmd[0]
	}

	args = append(args, config.Args...)

	vm.SetProcessAttr(&notto.ProcessAttr{
		Stdout:  config.Stdout,
		Stderr:  config.Stderr,
		Environ: notto.Environ(config.Env),
		Cwd:     config.WorkDir,
		Argv:    args,
	})

	return n, nil
}

func init() {
	exec.Register("notto", &NodeInterpreter{})
}
