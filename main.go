package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

// KE, KD and KI are the constant used for the PID computation
// MIN and MAX are min and max ouput values
// TARGET is the wanted angle between pod and check point, ie 0°
const (
	KE     float64 = 210.0
	KD     float64 = -0.02
	KI     float64 = -5
	MAX    float64 = 100.0
	MIN    float64 = 0.0
	TARGET float64 = 0.0
)

// Point implements a point in the plan.
type Point struct {
	x float64 // x coordinate
	y float64 // y coordinate
}

func getDistance(p1, p2 Point) float64 {
	var dist float64
	dist = math.Sqrt(math.Pow((p1.x-p2.x), 2) + math.Pow((p1.y-p2.y), 2))
	return dist
}

// PIDController implements a PID controller.
type PIDController struct {
	ke         float64   // proportional gain
	ki         float64   // integral gain
	kd         float64   // derivate gain
	setpoint   float64   // current setpoint
	integral   float64   // integral sum
	lastUpdate time.Time // time of last update
	prevError  float64   // previous error
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	var pc = PIDController{ke: KE, ki: KI, kd: KD, setpoint: TARGET, integral: 0, lastUpdate: time.Now(), prevError: 0}
	for {
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle float64
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)

		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)
		//dt := float64(time.Since(pc.lastUpdate))
		fmt.Fprintf(os.Stderr, "dt: %v\n", time.Since(pc.lastUpdate))
		fmt.Fprintf(os.Stderr, "angle: %v\n", nextCheckpointAngle)
		var dt = 0.1

		//		p1 := Point{x: x, y: y}
		//		p2 := Point{x: nextCheckpointX, y: nextCheckpointY}
		//		p3 := Point{x: opponentX, y: opponentY}

		error := math.Abs(math.Cos(nextCheckpointAngle - pc.setpoint))
		pc.integral = pc.integral + error*dt
		derivative := (error - pc.prevError) / dt
		var thrust = pc.ke*error + pc.ki*pc.integral + pc.kd*derivative
		if thrust > MAX {
			thrust = MAX
		} else if thrust < MIN {
			thrust = MIN
		}
		fmt.Fprintf(os.Stderr, "output: %f; %f; %f; %f\n", error, derivative, pc.integral, thrust)
		pc.prevError = error
		pc.lastUpdate = time.Now()

		// You have to output the target position
		// followed by the power (0 <= thrust <= 100)
		// i.e.: "x y thrust"
		fmt.Printf("%d %d %d\n", int(nextCheckpointX), int(nextCheckpointY), int(thrust))
	}
}
