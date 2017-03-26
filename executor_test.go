package exec

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestExecutor(t *testing.T) {

	conf := Config{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		//Interpreter: []string{"node", "-e"},
		Cmd: []string{"env"},
		//Script: "setTimeout(function () {console.log('helloe'); console.log(process.env)}, 1000)",
	}

	e := New(conf)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	if e := e.Start(ctx); e != nil {
		t.Fatal(e)
	}

}
