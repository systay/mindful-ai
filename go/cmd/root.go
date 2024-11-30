package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func Execute() {
	root := &cobra.Command{
		Use:   "mindful",
		Short: "Tool to create mindfulness meditations with",
	}

	root.CompletionOptions.HiddenDefaultCmd = true

	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func fail(err error, messages ...string) {
	if err == nil {
		return
	}
	context := strings.Join(messages, " ")

	if context != "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s - %s\n", context, err.Error())
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	os.Exit(1)
}
