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

var boostAvailable = false

// Point structure
type Point struct {
	x float64 // x coordinate
	y float64 // y coordinate
}

// Velocity structure
type Velocity struct {
	x float64 // x composante
	y float64 // y composante
}

// distance returns euclidian distance between two point in the plan
func distance(p1, p2 Point) float64 {
	var dist float64
	dist = math.Sqrt(math.Pow((p1.x-p2.x), 2) + math.Pow((p1.y-p2.y), 2))
	return dist
}

// getAngle gives the rotation angle of current point to aim a goal point
func getAngle(current, goal Point) float64 {
	var d = distance(current, goal)
	var dx = (goal.x - current.x) / d
	var dy = (goal.y - current.y) / d

	// radian to degre
	var a = math.Acos(dx) * 180.0 / math.Pi
	fmt.Fprintf(os.Stderr, "\tgetAngle: %f (before correction)\n", a)

	if dy < 0 {
		a = 360.0 - a
	}

	fmt.Fprintf(os.Stderr, "\tgetAngle: %f (after correction)\n", a)
	return a
}

// compute the angle between current direction and a goal point
func diffAngle(current, goal Point, angle float64) float64 {
	// get trigonometric angle in x,y plan
	var a = getAngle(current, goal)

	// compute angle if we turn in trigonometric way (turn right, because y axis goes down)
	var right, left float64

	if angle <= a {
		right = a - angle
	} else {
		right = 360.0 - angle + a
	}

	// compute angle if we turn in counter trigonometric way (turn left, because y axis goes down)
	if angle >= a {
		left = angle - a
	} else {
		left = angle + 360.0 - a
	}

	// choose the smaller angle
	if right < left {
		fmt.Fprintf(os.Stderr, "\tdiffAngle: right %f\n", right)
		return right
	}
	fmt.Fprintf(os.Stderr, "\tdiffAngle: left -%f\n", left)
	return -left
}

// rotate compute new angle (limited to a 18° variation)
// to rotate in wanted direction
func rotate(current, goal Point, angle float64) float64 {
	var a = diffAngle(current, goal, angle)
	var result float64
	// rotation is limited to 18° per tour
	if a > 18.0 {
		a = 18.0
	} else if a < -18.0 {
		a = -18.0
	}

	result = angle + a

	if result >= 360.0 {
		result = result - 360.0
	} else if result < 0.0 {
		result = result + 360.0
	}

	fmt.Fprintf(os.Stderr, "\trotate: %f\n", result)
	return result
}

// boost compute the new velocity from current velocity, thrust to apply
// and wanted angle
func boost(thrust, angle float64, v Velocity) Velocity {
	var vel Velocity
	/*    // N'oubliez pas qu'un pod qui a activé un shield ne peut pas accélérer pendant 3 tours
	      if (this.shield) {
	          return;
	      }*/

	// Conversion de l'angle en radian
	var ra = angle * math.Pi / 180.0

	// Trigonométrie
	vel.x = v.x + math.Cos(ra)*thrust
	vel.y = v.y + math.Sin(ra)*thrust

	fmt.Fprintf(os.Stderr, "\tvelocity: %f\n", vel)
	return vel
}

// move return new position according to current pos and desired velocity
func move(curPos Point, v Velocity, t float64) Point {
	var newPos Point
	newPos.x = curPos.x + v.x*t
	newPos.y = curPos.y + v.y*t

	fmt.Fprintf(os.Stderr, "\tmove to newPos: (%f, %f)\n", newPos.x, newPos.y)
	return newPos
}

// play predict new position and velocity at the end of the tour based:
// - current position, goal position
// - current angle and velocity
// - and thrust to apply
func play(current, goal Point, angle, thrust float64, v Velocity) (Point, Velocity) {
	var a = rotate(current, goal, angle)
	var vel = boost(thrust, a, v)
	var p = move(current, vel, 1.0)
	vel.x = vel.x * 0.85
	vel.y = vel.y * 0.85

	fmt.Fprintf(os.Stderr, "\tnewPos: (%f, %f) newVel: (%f, %f)\n", p.x, p.y, vel.x, vel.y)
	return p, vel
}

// getThrust
func getThrust(angle, dist float64) (float64, bool) {
	var output = 100.0
	var boost = false
	angle = math.Abs(angle)

	if boostAvailable && dist > boostRadius && angle < 4 {
		boost = true
	}

	if angle > 90 {
		output = 0
	} else if angle > 15 {
		output = 55
	} else if dist <= brakeStep3 {
		output = 25
	} else if dist <= brakeStep2 {
		output = 50
	} else if dist <= brakeStep1 {
		output = 75
	}

	fmt.Fprintf(os.Stderr, "\tangle: %f; dist: %f; out: %f\n", angle, dist, output)
	return output, boost
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	var v Velocity
	var p Point
	var loop = 1
	for {
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle float64
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)
		fmt.Fprintf(os.Stderr, "### T O U R   %d   ###\n", loop*2)
		fmt.Fprintf(os.Stderr, "x: %f; y: %f, d: %f, a: %f\n", x, y, nextCheckpointDist, nextCheckpointAngle)

		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)

		thrust, boost := getThrust(nextCheckpointAngle, nextCheckpointDist)
		var angle = getAngle(Point{x, y}, Point{nextCheckpointX, nextCheckpointY})
		p = Point{x, y}
		for i := loop + 1; i < loop+2; i++ {
			fmt.Fprintf(os.Stderr, "- simulation loop %d   ###\n", i*2)
			p, v = play(p, Point{nextCheckpointX, nextCheckpointY}, angle+nextCheckpointAngle, thrust, v)
		}

		loop++
		// You have to output the target position
		// followed by the power (0 <= thrust <= 100)
		// i.e.: "x y thrust"
		if boost == true && boostAvailable == true {
			fmt.Printf("%d %d BOOST BOOST\n", int(nextCheckpointX), int(nextCheckpointY))
			boostAvailable = false
		} else {
			fmt.Printf("%d %d %d (%d,%d°,%d)\n", int(nextCheckpointX), int(nextCheckpointY), int(thrust), int(thrust), int(nextCheckpointAngle), int(nextCheckpointDist))
		}
	}
}
