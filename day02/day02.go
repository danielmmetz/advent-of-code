package day02

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/danielmmetz/advent-of-code/errors"
)

var Cmd = cobra.Command{
	Use:   "2",
	Short: "day 2",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runE()
	},
}

func init() {
	Cmd.Flags().IntVar(&part, "part", 1, "part two")
}

var part int

func runE() error {
	if err := validate(); err != nil {
		return err
	}

	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("error reading from stdin: %v", err)
	}

	inputText := strings.TrimSpace(string(bs))
	rawValues := strings.Split(inputText, ",")
	ints := make([]int, 0, len(rawValues))
	for _, c := range rawValues {
		i, err := strconv.Atoi(c)
		if err != nil {
			return fmt.Errorf("error constructing intcodeProgram: cannot convert to int: %s", c)
		}
		ints = append(ints, i)
	}

	switch part {
	case 1:
		return partOne(ints)
	case 2:
		return partTwo(ints)
	default:
		return fmt.Errorf("invalid part specified: %d", part)
	}
}

func partOne(ints []int) error {
	ints[1] = 12
	ints[2] = 2
	program := intcodeProgram(ints)
	if err := program.run(); err != nil {
		return fmt.Errorf("error while running intcode program: %v", err)
	}

	fmt.Println(program[0])
	return nil
}

func partTwo(ints []int) error {
	const targetOutput = 19690720
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			dupe := make([]int, len(ints))
			copy(dupe, ints)
			dupe[1], dupe[2] = i, j
			program := intcodeProgram(dupe)
			if err := program.run(); err != nil {
				return fmt.Errorf("error while running intcode program: %v", err)
			}
			if program[0] == targetOutput {
				fmt.Println(100*i + j)
				return nil
			}
		}
	}
	return fmt.Errorf("no noun-verb combo found")
}

type intcodeProgram []int

func (p intcodeProgram) run() error {
	var err error
	for i := 0; err == nil; i += 4 {
		err = operationFrom(p, i).apply(p)
	}

	if err != errHalt {
		return err
	}

	return nil
}

func (p intcodeProgram) equal(o intcodeProgram) bool {
	if len(p) != len(o) {
		return false
	}
	for i := range p {
		if p[i] != o[i] {
			return false
		}
	}
	return true
}

type opcode int

const (
	add      opcode = 1
	multiply opcode = 2
	halt     opcode = 99
)

var errHalt error = fmt.Errorf("halt")

type operation struct {
	opcode opcode
	arg1   int
	arg2   int
	output int
}

func operationFrom(p intcodeProgram, index int) operation {
	if opcode(p[index]) == halt {
		return operation{opcode: halt}
	}
	return operation{
		opcode: opcode(p[index]),
		arg1:   p[index+1],
		arg2:   p[index+2],
		output: p[index+3],
	}
}

func (o operation) apply(p intcodeProgram) error {
	switch o.opcode {
	case halt:
		return errHalt
	case add:
		p[o.output] = p[o.arg1] + p[o.arg2]
	case multiply:
		p[o.output] = p[o.arg1] * p[o.arg2]
	default:
		return fmt.Errorf("unknown opcode: %d", int(o.opcode))
	}
	return nil
}

func validate() error {
	cases := []struct {
		input    intcodeProgram
		expected intcodeProgram
	}{
		{intcodeProgram([]int{1, 0, 0, 0, 99}), intcodeProgram([]int{2, 0, 0, 0, 99})},
		{intcodeProgram([]int{2, 3, 0, 3, 99}), intcodeProgram([]int{2, 3, 0, 6, 99})},
		{intcodeProgram([]int{2, 4, 4, 5, 99, 0}), intcodeProgram([]int{2, 4, 4, 5, 99, 9801})},
		{intcodeProgram([]int{1, 1, 1, 4, 99, 5, 6, 0, 99}), intcodeProgram([]int{30, 1, 1, 4, 2, 5, 6, 0, 99})},
	}

	var results errors.TestResults
	for i, c := range cases {
		err := c.input.run()
		if err != nil {
			results.AppendFailure(fmt.Sprintf("unexpected error in test case %d: %v", i, err))
			continue
		}
		if !c.input.equal(c.expected) {
			results.AppendFailure(fmt.Sprintf("expected test case %d to equal %v, got %v", i, c.expected, c.input))
		}
	}

	return results.Err()
}
