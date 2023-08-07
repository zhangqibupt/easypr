package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func Test(t *testing.T) {
	CreateRun()(&cobra.Command{}, []string{})
}
