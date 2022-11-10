package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := cmd[0]
	args := cmd[1:]
	c := exec.Command(command, args...)

	envVars := []string{}
	for k, v := range env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v.Value))
	}
	c.Env = append(os.Environ(), envVars...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		log.Println(err)
	}

	return c.ProcessState.ExitCode()
}
