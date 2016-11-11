package exec

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestExecutor(t *testing.T) {

	conf := Config{
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Stdin:       os.Stdin,
		Interpreter: []string{"node", "-e"},
		//Cmd:         []string{"echo", "hello, world"},
		Script: "setTimeout(function () {console.log('helloe')}, 10000)",
	}

	e := New(conf)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	if e := e.Start(ctx); e != nil {
		t.Fatal(e)
	}

}
