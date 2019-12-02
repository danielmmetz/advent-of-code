package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if errs := validate(); len(errs) != 0 {
		fmt.Fprintf(os.Stderr, "failed the following test cases:\n")
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "\t%v\n", err)
		}
		os.Exit(1)
	}

    holistic := flag.Bool("holisitic", false, "include accounting for fuel usage")
    flag.Parse()

	var total int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		mass, err := strconv.Atoi(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "expected int, got %s\n", err)
		}
		total += FuelRequired(*holistic, mass)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(total)
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

func validate() []error {
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

	var errors []error
	for _, c := range cases {
		actual := FuelRequired(c.holisitc, c.mass)
		if actual != c.expected {
			errors = append(errors, fmt.Errorf("expected FuelRequired(%t, %d) == %d, got %d", c.holisitc, c.mass, c.expected, actual))
		}
	}

	return errors
}
