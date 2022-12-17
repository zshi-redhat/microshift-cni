package util

import (
	"bytes"
	"fmt"
	"strings"

	"k8s.io/klog/v2"
	kexec "k8s.io/utils/exec"
)

var runner kexec.Interface

func init() {
	runner = kexec.New()
}

func RunCommand(bin string, args ...string) error {
	path, err := runner.LookPath(bin)
	if err != nil {
		return err
	}
	var stdout, stderr bytes.Buffer
	cmd := runner.Command(path, args...)
	cmd.SetStdout(&stdout)
	cmd.SetStderr(&stderr)

	argStr := strings.Join(args, " ")
	klog.V(5).Infof("Exec: %s %s", path, argStr)

	err = cmd.Run()
	if err != nil {
		stderrStr := stderr.String()
		klog.Errorf("Exec %s %s : stderr: %q", path, argStr, stderrStr)
		return fmt.Errorf("failed to run '%s %s': %v\n  %q", path, argStr, err, stderrStr)
	}
	stdoutStr := stdout.String()
	klog.Infof("Exec: %s %s: stdout: %q", path, argStr, stdoutStr)
	return nil
}
