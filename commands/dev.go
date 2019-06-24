package commands

import (
	"github.com/spf13/cobra"
)

// NewDevelopCommand create dev command
func NewDevelopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "command set just for develop",
	}
	return cmd
}
