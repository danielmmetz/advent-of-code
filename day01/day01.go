package day01

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/danielmmetz/advent-of-code/errors"
)

var Cmd = cobra.Command{
	Use:   "1",
	Short: "day 1",
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

	flag.Parse()

	var total int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		mass, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("expected int, got %s", err)
		}
		switch part {
		case 1:
			total += FuelRequired(false, mass)
		case 2:
			total += FuelRequired(true, mass)
		default:
			return fmt.Errorf("invalid part specified: %d", part)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println(total)
	return nil
}

func FuelRequired(holisitc bool, mass int) int {
	if holisitc {
		return holisticFuelRequired(mass)
	}
	return fuelRequired(mass)
}

func fuelRequired(mass int) int {
	return mass/3 - 2
}

func holisticFuelRequired(mass int) int {
	var total int
	for {
		required := fuelRequired(mass)
		if required <= 0 {
			return total
		}
		total += required
		mass = required
	}
}

func validate() error {
	cases := []struct {
		holisitc       bool
		mass, expected int
	}{
		{false, 12, 2},
		{false, 14, 2},
		{false, 1969, 654},
		{false, 100756, 33583},
		{true, 14, 2},
		{true, 1969, 966},
		{true, 100756, 50346},
	}

	var results errors.TestResults
	for _, c := range cases {
		actual := FuelRequired(c.holisitc, c.mass)
		if actual != c.expected {
			results.AppendFailure(fmt.Sprintf("expected FuelRequired(%t, %d) == %d, got %d", c.holisitc, c.mass, c.expected, actual))
		}
	}

	return results.Err()
}
