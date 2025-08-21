package util

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiBlue   = "\033[34m"
	ansiCyan   = "\033[36m"
)

type ContainerLogger struct {
	testing.TB
	containerName string
	colorPrefix   string
}

type TcLogger struct {
	testing.TB
}

func (t TcLogger) Printf(format string, v ...any) {
	t.Helper()
	prefix := fmt.Sprintf("%s[%s]%s", ansiCyan, "testcontainers", ansiReset)
	line := fmt.Sprintf(format, v...)
	fmt.Fprintf(os.Stdout, "%s %s\n", prefix, line)
}

func (c *ContainerLogger) Accept(log testcontainers.Log) {
	content := string(log.Content)
	lines := strings.Split(content, "\n")

	prefix := fmt.Sprintf("%s[%s]%s", c.colorPrefix, c.containerName, ansiReset)
	for i, line := range lines {
		// Skip final empty slice element caused by trailing newline to avoid double blank line
		if i == len(lines)-1 && line == "" {
			continue
		}

		if strings.ToUpper(log.LogType) == "STDERR" {
			fmt.Fprintf(os.Stdout, "%s %s%s%s\n", prefix, ansiBold+ansiRed, line, ansiReset)
		} else {
			fmt.Fprintf(os.Stdout, "%s %s\n", prefix, line)
		}
	}
}
