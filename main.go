package main

import (
	"fmt"
	"math"
	"os"
	//	"time"
)

// Minimum distance for activiating boost
const boostRadius = 2000

// Distance starting from the middle of the checkpoint for the racer to aim for
const radius = 350

// Distance steps for slowing down the racer
const brakeStep1 = 1300
const brakeStep2 = 1100
const brakeStep3 = 800

var boostAvailable = true

// Point structure
type Point struct {
	x float64 // x coordinate
	y float64 // y coordinate
}

func getDistance(p1, p2 Point) float64 {
	var dist float64
	dist = math.Sqrt(math.Pow((p1.x-p2.x), 2) + math.Pow((p1.y-p2.y), 2))
	return dist
}

// FeedForward ...
func getThrust(angle, dist float64) (float64, bool) {
	var output = 100.0
	var boost = false
	angle = math.Abs(angle)

	if boostAvailable && dist > boostRadius && angle < 4 {
		boost = true
	}

	if angle > 90 {
		output = 0
	} else if dist <= brakeStep3 {
		output = 25
	} else if dist <= brakeStep2 {
		output = 50
	} else if dist <= brakeStep1 {
		output = 75
	}

	fmt.Fprintf(os.Stderr, "angle: %f; dist: %f; out: %f\n", angle, dist, output)
	return output, boost
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	for {
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle float64
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)

		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)

		thrust, boost := getThrust(nextCheckpointAngle, nextCheckpointDist)

		// You have to output the target position
		// followed by the power (0 <= thrust <= 100)
		// i.e.: "x y thrust"
		if boost == true && boostAvailable == true {
			fmt.Printf("%d %d BOOST BOOST\n", int(nextCheckpointX), int(nextCheckpointY))
			boostAvailable = false
		} else {
			fmt.Printf("%d %d %d %d\n", int(nextCheckpointX), int(nextCheckpointY), int(thrust), int(thrust))
		}
	}
}
