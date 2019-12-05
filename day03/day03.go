package day03

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/danielmmetz/advent-of-code/errors"
)

var Cmd = cobra.Command{
	Use:   "3",
	Short: "day 3",
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

	paths := make([]path, 0, 2)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		path, err := newPath(input)
		if err != nil {
			return fmt.Errorf("error constructing path: %v", err)
		}
		paths = append(paths, path)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if len(paths) != 2 {
		return fmt.Errorf("expected 2 paths, got %d", len(paths))
	}

    var answer int

	switch part {
	case 1:
        answer = distanceOfClosestIntersectionByManhattan(paths[0], paths[1])
	case 2:
        answer = distanceOfClosestIntersectionByDelay(paths[0], paths[1])
	default:
		return fmt.Errorf("invalid part specified: %d", part)
	}

    fmt.Println(answer)
	return nil
}

type point struct {
	x, y int
}

type points []point

func (p points) Map() map[point]bool {
    m := make(map[point]bool)
    for _, p := range p {
        m[p] = true;
    }
    return m
}

// pointsEncountered returns a temporally-sorted list of integral-valued points
// encountered by processing a pathMove from start.
func pointsEncountered(start point, move pathMove) points {
	var xMultiplier, yMultiplier int

	switch move.direction {
	case up:
		yMultiplier = 1
	case down:
		yMultiplier = -1
	case left:
		xMultiplier = -1
	case right:
		xMultiplier = 1
	}

	current := start
	points := []point{}
	for i := 0; i < move.magnitutde; i++ {
		current.x += xMultiplier
		current.y += yMultiplier
		points = append(points, current)
	}
	return points
}

func manhattanDistance(a, b point) int {
	return int(math.Abs(float64(a.x-b.x)) + math.Abs(float64(a.y-b.y)))
}

type pathMove struct {
	direction  direction
	magnitutde int
}

type direction string

const (
	up    direction = "U"
	down  direction = "D"
	left  direction = "L"
	right direction = "R"
)

// newPathMove returns a pathmove parsed from input.
// input is expected to look like "R5" or "D2""
func newPathMove(input string) (pathMove, error) {
	if input == "" {
		return pathMove{}, fmt.Errorf("empty input")
	}

	var pm pathMove
	switch direction(input[0]) {
	case up, down, left, right:
		pm.direction = direction(input[0])
	default:
		return pm, fmt.Errorf("invalid direction: %s", string(input[0]))
	}

	magnitude, err := strconv.Atoi(input[1:])
	if err != nil {
		return pm, fmt.Errorf("invalid magnitude: %s", input[1:])
	}
	pm.magnitutde = magnitude
	return pm, nil
}

type path []pathMove

// newPath returns a path parsed from input.
// input is expected to look like "R5,D2,..."
func newPath(input string) (path, error) {
	path := []pathMove{}

	steps := strings.Split(input, ",")
	for _, step := range steps {
		pm, err := newPathMove(step)
		if err != nil {
			return path, err
		}
		path = append(path, pm)
	}
	return path, nil
}

func collectPointsAlong(p path) points {
    points := []point{}

	var current point
	for _, step := range p {
		encounteredPoints := pointsEncountered(current, step)
		current = encounteredPoints[len(encounteredPoints)-1]
		for _, point := range encounteredPoints {
			points = append(points, point)
		}
	}
	return points
}

func intersection(a, b map[point]bool) points {
	if len(a) > len(b) {
		b, a = a, b
	}

	common := []point{}
	for point := range a {
		if b[point] {
			common = append(common, point)
		}
	}

	return common
}

func delayToFirstIntersection(a, b points) int {
    intersections := intersection(a.Map(), b.Map())

    path1Delay := make([]int, 0, len(intersections))
    path2Delay := make([]int, 0, len(intersections))

    for _, p := range intersections {
        for i, p1 := range a {
            if p == p1 {
                path1Delay = append(path1Delay, i+1)
                break
            }
        }
        for i, p2 := range b {
            if p == p2 {
                path2Delay = append(path2Delay, i+1)
                break
            }
        }
    }

    minCombo := -1
    for i := range(intersections) {
        delay := path1Delay[i] + path2Delay[i]
        if minCombo < 0 || delay < minCombo {
            minCombo = delay
        }
    }
    return minCombo
}

func closestToOrigin(candidates points) point {
	var winner point
	minDistanceSoFar := -1
	for _, p := range candidates {
		distance := manhattanDistance(p, point{x: 0, y: 0})
		if minDistanceSoFar < 0 || distance < minDistanceSoFar {
			winner = p
			minDistanceSoFar = distance
		}
	}

	return winner
}

func closestIntersectionToOrigin(p1, p2 path) point {
	p1points := collectPointsAlong(p1)
	p2points := collectPointsAlong(p2)
	intersectionPoints := intersection(p1points.Map(), p2points.Map())
	return closestToOrigin(intersectionPoints)
}

func distanceOfClosestIntersectionByManhattan(p1, p2 path) int {
	return manhattanDistance(closestIntersectionToOrigin(p1, p2), point{x: 0, y: 0})
}

func distanceOfClosestIntersectionByDelay(p1, p2 path) int {
	p1points := collectPointsAlong(p1)
	p2points := collectPointsAlong(p2)
    return delayToFirstIntersection(p1points, p2points)
}

func validate() error {
	cases := []struct {
		path1    string
		path2    string
        distanceFunc func(p1, p2 path) int
		expected int
	}{
		{
            "R8,U5,L5,D3",
            "U7,R6,D4,L4",
            distanceOfClosestIntersectionByManhattan,
            6,
        },
		{
            "R75,D30,R83,U83,L12,D49,R71,U7,L72",
            "U62,R66,U55,R34,D71,R55,D58,R83",
            distanceOfClosestIntersectionByManhattan,
            159,
        },
		{
            "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51",
            "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
            distanceOfClosestIntersectionByManhattan,
            135,
        },
		{
            "R8,U5,L5,D3",
            "U7,R6,D4,L4",
            distanceOfClosestIntersectionByDelay,
            30,
        },
		{
            "R75,D30,R83,U83,L12,D49,R71,U7,L72",
            "U62,R66,U55,R34,D71,R55,D58,R83",
            distanceOfClosestIntersectionByDelay,
            610,
        },
		{
            "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51",
            "U98,R91,D20,R16,D67,R40,U7,R15,U6,R7",
            distanceOfClosestIntersectionByDelay,
            410,
        },
	}

	var results errors.TestResults
	for i, c := range cases {
		path1, err1 := newPath(c.path1)
		path2, err2 := newPath(c.path2)
		if err1 != nil {
			results.AppendFailure(fmt.Sprintf("unexpected error in test case %d converting path1: %v", i, err1))
			continue
		}
		if err2 != nil {
			results.AppendFailure(fmt.Sprintf("unexpected error in test case %d converting path2: %v", i, err2))
			continue
		}
		actual := c.distanceFunc(path1, path2)
		if actual != c.expected {
			results.AppendFailure(fmt.Sprintf("expected test case %d to equal %v, got %v", i, c.expected, actual))
		}
	}

	return results.Err()
}
