package main

import (
	"os"
	"os/exec"
	"strings"
)

// the string inside this will actually be excuted as a command on the system
var c = `
START EXEC OMIT
go run ./cmd/jwt-issue auth.ed
END EXEC OMIT
`

func main() {
	c = strings.TrimSpace(c)

	// go present actually removes the OMIT lines
	// before this gets executed so it may not be here
	// so only try to split if we see newlines
	if strings.Contains(c, "\n") {
		parts := strings.Split(c, "\n")
		c = parts[1]
	}
	cmdParts := strings.Fields(c)

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
