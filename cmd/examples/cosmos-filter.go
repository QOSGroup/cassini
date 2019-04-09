// +build ignore

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var logFile string

func main() {
	root := newRootCommand()
	root.AddCommand(newParseCommand())

	if err := root.Execute(); err != nil {
		fmt.Println("Exit by error: ", err)
	}
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "address-filter",
		Short: "filter address in log file",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return
		},
	}
	return root
}

func newParseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "parse address from log file",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Parsing log: ", logFile)
			parseFile(logFile)
			fmt.Println("done.")
			return nil
		},
	}

	addFlags(cmd)
	return cmd
}

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&logFile,
		"log", "./easyzone.log", "log file path")
}

func parseFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("file open error: ", file)
		return
	}
	buf := bufio.NewReader(f)
	var line string
	line, err = buf.ReadString('\n')
	for err == nil {
		line = strings.TrimSpace(line)
		strs := strings.Split(line, "------")
		if len(strs) > 2 {
			ss := strings.Split(strs[1], "\\x22")
			if len(ss) > 2 {
				fmt.Println(ss[3])
			}
		}
		line, err = buf.ReadString('\n')
	}
	if err != nil && err != io.EOF {
		fmt.Println("parse file error: ", err)
		return
	}
}
