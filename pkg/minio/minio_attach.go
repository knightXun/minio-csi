package minio

import (
	"k8s.io/klog"
	"os"
	"os/exec"
)

const (
	envHostRootFS = "HOST_ROOTFS"
	s3fs      = "s3fs"
)

var (
	hostRootFS = "/"
	hasNBD     = false
)

func init() {
	host := os.Getenv(envHostRootFS)
	if len(host) > 0 {
		hostRootFS = host
	}
}

// Check if s3fs tools are installed.
func checks3fsTools() bool {
	_, err := execCommand("/usr/bin/s3fs", []string{"--help"})
	if err != nil {
		klog.V(3).Infof("s3fs tools not found")
		return false
	}
	klog.V(3).Infof("s3fs tools were found.")
	return true
}

func execCommand(command string, args []string) ([]byte, error) {
	// #nosec
	cmd := exec.Command(command, args...)
	return cmd.CombinedOutput()
}
