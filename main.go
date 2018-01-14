package main

import (
	"fmt"
	"math"
	"os"
	//	"time"
)

// FeedForward ...
func FeedForward(thrust, angle, dist float64, boosted bool) (float64, bool) {
	var output = thrust
	var boost = false
	angle = math.Abs(angle)
	fmt.Fprintf(os.Stderr, "thrust: %f; angle: %f; dist: %f; out: %f\n", thrust, angle, dist, output)

	if !boosted && dist > 6500 && angle < 4 {
		boost = true
	}

	if dist > 2500 && angle < 10 {
		output = 100
	} else if dist < 1150 && angle < 10 {
		output = 0
	} else if angle > 60 {
		output = 50
	} else {
		output = 85
	}
	return output, boost
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	var boosted = false
	for {
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle float64
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)

		var opponentX, opponentY int
		var boost = false
		fmt.Scan(&opponentX, &opponentY)
		thrust := 100.0
		thrust, boost = FeedForward(thrust, nextCheckpointAngle, nextCheckpointDist, boosted)

		// You have to output the target position
		// followed by the power (0 <= thrust <= 100)
		// i.e.: "x y thrust"
		if boost == true && !boosted {
			fmt.Printf("%d %d BOOST BOOST\n", int(nextCheckpointX), int(nextCheckpointY))
			boosted = true
		} else {
			fmt.Printf("%d %d %d %d\n", int(nextCheckpointX), int(nextCheckpointY), int(thrust), int(thrust))
		}
	}
}
