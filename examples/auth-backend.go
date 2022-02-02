package main

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// the string inside this will actually be excuted as a command on the system
var commands = `
START EXEC OMIT
go run ./cmd/2-backend auth.ed.pub & sleep 2
go run ./cmd/grpc-local
END EXEC OMIT
`

func main() {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	scriptPath := filepath.Join(dir, "script.sh")

	err = ioutil.WriteFile(scriptPath, []byte(commands), 0700)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// set child process group id
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	kill := func() {
		syscall.Kill(-cmd.Process.Pid, 15)
	}
	defer kill()

	cmd.Process.Wait()

	// second second kill as a lazy way to be sure we don't exit too soon
	// this is ugly but other methods don't seem to work
	// reliably when running through go present
	kill()
}
