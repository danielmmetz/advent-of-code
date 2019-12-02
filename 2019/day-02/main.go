package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	if errs := validate(); len(errs) != 0 {
		fmt.Fprintln(os.Stderr, "failed the following test cases:")
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, "\t", err)
		}
		os.Exit(1)
	}

	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin: ", err)
		os.Exit(1)
	}

	inputText := strings.TrimSpace(string(bs))
	rawValues := strings.Split(inputText, ",")
	ints := make([]int, 0, len(rawValues))
	for _, c := range rawValues {
		i, err := strconv.Atoi(c)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error constructing intcodeProgram: cannot convert to int: ", c)
			os.Exit(1)
		}
		ints = append(ints, i)
	}

	ints[1] = 12
	ints[2] = 2
	program := intcodeProgram(ints)
	if err := program.run(); err != nil {
		fmt.Fprintf(os.Stderr, "error while running intcode program: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(program[0])
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

func validate() []error {
	cases := []struct {
		input    intcodeProgram
		expected intcodeProgram
	}{
		{intcodeProgram([]int{1, 0, 0, 0, 99}), intcodeProgram([]int{2, 0, 0, 0, 99})},
		{intcodeProgram([]int{2, 3, 0, 3, 99}), intcodeProgram([]int{2, 3, 0, 6, 99})},
		{intcodeProgram([]int{2, 4, 4, 5, 99, 0}), intcodeProgram([]int{2, 4, 4, 5, 99, 9801})},
		{intcodeProgram([]int{1, 1, 1, 4, 99, 5, 6, 0, 99}), intcodeProgram([]int{30, 1, 1, 4, 2, 5, 6, 0, 99})},
	}

	var errors []error
	for i, c := range cases {
		err := c.input.run()
		if err != nil {
			errors = append(errors, fmt.Errorf("unexpected error in test case %d: %w", i, err))
			continue
		}
		if !c.input.equal(c.expected) {
			errors = append(errors, fmt.Errorf("expected test case %d to equal %v, got %v", i, c.expected, c.input))
		}
	}

	return errors
}
