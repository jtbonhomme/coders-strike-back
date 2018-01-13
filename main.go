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
	KE     float64 = 0.0025
	KD     float64 = -0.025
	KI     float64 = 0.10
	MAX    float64 = 100.0
	MIN    float64 = 5.0
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

// GetPID compute PID
func GetPID(pc *PIDController, value float64, dt float64) float64 {
	error := value - pc.setpoint
	pc.integral = pc.integral + error*dt
	derivative := (error - pc.prevError) / dt
	var output = pc.ke*error + pc.ki*pc.integral + pc.kd*derivative
	if output > MAX {
		output = MAX
	} else if output < MIN {
		output = MIN
	}
	pc.prevError = error
	fmt.Fprintf(os.Stderr, "ke*e: %f; kd*de: %f; ki*ie: %f; o: %f\n",  pc.ke*error, pc.kd*derivative, pc.ki*pc.integral, output)
	fmt.Fprintf(os.Stderr, "err : %f; deriv: %f; integ: %f; o: %f\n", error, derivative, pc.integral, output)
	return output
}

// FeedForward
func FeedForward(thrust, angle, dist float64) float64 {
    var output = thrust
    fmt.Fprintf(os.Stderr, "thrust: %f; angle: %f; dist: %f; out: %f\n", thrust, angle, dist, output)
    if angle == 0 && dist > 1500 {
        output = 100
    } else if angle == 0 && dist < 1500 {
        output = 60
    }
	return output
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	var pc = PIDController{ke: KE, ki: KI, kd: KD, setpoint: TARGET, integral: 0, lastUpdate: time.Now(), prevError: 20}
	for {
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle float64
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)

		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)
		dt:= float64(time.Since(pc.lastUpdate))/1e9
		fmt.Fprintf(os.Stderr, "dt: %f angle: %v dist : %v\n", dt, nextCheckpointAngle, nextCheckpointDist)

		//		p1 := Point{x: x, y: y}
		//		p2 := Point{x: nextCheckpointX, y: nextCheckpointY}
		//		p3 := Point{x: opponentX, y: opponentY}

		thrust := GetPID(&pc, nextCheckpointDist, dt)
		pc.lastUpdate = time.Now()
        thrust = FeedForward(thrust, nextCheckpointAngle, nextCheckpointDist)
        
		// You have to output the target position
		// followed by the power (0 <= thrust <= 100)
		// i.e.: "x y thrust"
		fmt.Printf("%d %d %d\n", int(nextCheckpointX), int(nextCheckpointY), int(thrust))
	}
}
