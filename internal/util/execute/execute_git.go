package execute

import (
	"bytes"
	"os/exec"
)

func ExecuteGit(cmd *exec.Cmd, directory string) (string, error) {
	cmd.Dir = directory
	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	err := cmd.Run()
	return outputBuffer.String(), err
}
