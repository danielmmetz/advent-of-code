package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/danielmmetz/advent-of-code/day01"
	"github.com/danielmmetz/advent-of-code/day02"
	"github.com/danielmmetz/advent-of-code/day03"
)

var cmd = cobra.Command{
	Use: "aoc <day>",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf(`error: missing required argument "day"`)
	},
}

func init() {
	cmd.AddCommand(&day01.Cmd)
	cmd.AddCommand(&day02.Cmd)
	cmd.AddCommand(&day03.Cmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}
