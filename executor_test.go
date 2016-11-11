package exec

import (
	"os"
	"testing"
)

func TestExecutor(t *testing.T) {

	conf := Config{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		//Interpreter: []string{"sh", "-c"},
		Cmd: []string{"sleep", "10"},
	}

	e := New(conf)

	if e := e.Start(); e != nil {
		t.Fatal(e)
	}

}
