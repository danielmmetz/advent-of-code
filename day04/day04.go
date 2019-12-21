package day04

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"

	"github.com/danielmmetz/advent-of-code/errors"
)

var Cmd = cobra.Command{
	Use:   "4",
	Short: "day 4",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runE()
	},
}

func init() {
	Cmd.Flags().IntVar(&part, "part", 1, "part two")
}

var part int

const passwordLength = 6

var (
	constraintsPart1 = []func([]int) bool{
		hasTwoAdjacentDigits,
		isMonotonicallyNonDecreasing,
	}
	constraintsPart2 = []func([]int) bool{
		hasTwoAdjacentDigits,
		isMonotonicallyNonDecreasing,
		hasIndependentPair,
	}
)

func runE() error {
	if err := validate(); err != nil {
		return err
	}

	var constraints []func([]int) bool
	switch part {
	case 1:
		constraints = constraintsPart1
	case 2:
		constraints = constraintsPart2
	default:
		return fmt.Errorf("invalid part specified: %d", part)
	}

	var count int
	for i := 124075; i < 580769; i++ {
		if ValidPassword(i, constraints) {
			count++
		}
	}
	fmt.Println(count)
	return nil
}

func ValidPassword(candidate int, constraints []func([]int) bool) bool {
	digits := make([]int, passwordLength, passwordLength)
	if candidate >= int(math.Pow10(passwordLength)) {
		return false
	}
	for i := 0; i < passwordLength; i++ {
		digits[passwordLength-i-1] = candidate % 10
		candidate = candidate / 10
	}
	for _, constraint := range constraints {
		if !constraint(digits) {
			return false
		}
	}
	return true
}

func hasTwoAdjacentDigits(digits []int) bool {
	if len(digits) < 2 {
		return false
	}
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] == digits[i+1] {
			return true
		}
	}
	return false
}

func isMonotonicallyNonDecreasing(digits []int) bool {
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] > digits[i+1] {
			return false
		}
	}
	return true
}

func hasIndependentPair(digits []int) bool {
	repeatCount := 1
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] == digits[i+1] {
			repeatCount++
			continue
		}
		if repeatCount == 2 {
			return true
		}
		repeatCount = 1
	}
	return repeatCount == 2
}

func validate() error {
	cases := []struct {
		input       int
		constraints []func([]int) bool
		valid       bool
	}{
		{111111, constraintsPart1, true},
		{223450, constraintsPart1, false},
		{123789, constraintsPart1, false},
		{112233, constraintsPart2, true},
		{123444, constraintsPart2, false},
		{111122, constraintsPart2, true},
	}

	var results errors.TestResults
	for i, c := range cases {
		actual := ValidPassword(c.input, c.constraints)
		if actual != c.valid {
			results.AppendFailure(fmt.Sprintf("expected test case %d to equal %v, got %v", i, c.valid, actual))
		}
	}

	return results.Err()
}
