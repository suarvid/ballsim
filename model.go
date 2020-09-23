package main

import (
	"fmt"
	"math"
)

type Model struct {
	width, height float64
	balls         []Ball
}

func createModel(width, height float64, balls []Ball) *Model {
	m := new(Model)
	m.width = width
	m.height = height
	m.balls = balls
	return m
}

// Used to have phi and r as fields in the struct
// but ended up not using them, they made it confusing as hell
type Ball struct {
	x, y, vx, vy, bisX, bisY float64
	mass, radius             float64
}

func createBall(x, y, vx, vy float64, mass, radius float64) *Ball {
	ball := new(Ball)
	ball.x = x
	ball.y = y
	ball.vx = vx
	ball.vy = vy
	ball.mass = mass
	ball.radius = radius
	return ball
}

// If we use range here we don't change the values of the actual balls
// It sort of becomes pass by value for some reason?
func (m Model) step(deltaT float64) {
	for i := 0; i < len(m.balls); i++ {
		m.edgeDetection(&m.balls[i], deltaT)
		m.moveBall(&m.balls[i], deltaT)
		m.applyGravity(&m.balls[i], 0.0)
		m.changeVelocity(&m.balls[i])
	}
	m.detectCollision(deltaT)
}

func (m Model) moveBall(b *Ball, deltaT float64) {
	b.x += b.vx * deltaT
	b.y += b.vy * deltaT
}

func (m Model) changeVelocity(b *Ball) {
	b.vx += b.bisX
	b.vy += b.bisY
}

func (m Model) applyGravity(b *Ball, g float64) {
	b.bisY = -g
}

// A bit unclear why the balls get to move a little outside the window before bouncing
// Could be something iffy with how the window is actually drawn
func (m Model) edgeDetection(b *Ball, deltaT float64) {
	if (b.x+b.vx*deltaT) < 0 || (b.x+b.vx*deltaT) > m.width {
		b.vx *= -1
	}
	// Should this really be b.radius? Maybe!
	if (b.y) < 0 || (b.y+b.vy*deltaT) > m.height {

		b.vy *= -1
	}
}

// This one gets its own nested loops, looked so messy otherwise
func (m Model) detectCollision(deltaT float64) {
	for i := 0; i < len(m.balls)-1; i++ {
		for j := i + 1; j < len(m.balls); j++ {
			if calculateDistance(m.balls[i], m.balls[j])+getVelocity(m.balls[i])*deltaT+getVelocity(m.balls[j])*deltaT < (m.balls[i].radius) {
				fmt.Println("Collision Detected")
				handleCollision(&m.balls[i], &m.balls[j], deltaT)
			}
		}
	}
}

// Returns the distance between two balls
func calculateDistance(b1, b2 Ball) float64 {
	distX := math.Abs(b1.x - b2.x)
	distY := math.Abs(b1.y - b2.y)
	return math.Sqrt(math.Pow(distX, 2) + math.Pow(distY, 2))
}

func handleCollision(b1, b2 *Ball, deltaT float64) {

	// Pretty sure these have to be b2 - b1, otherwise angles get iffy later
	xDiff := b2.x - b1.x
	yDiff := b2.y - b1.y

	/*
		When I write "the x-axis line" below, I mean a line perpendicular to the x-axis
		that passes through the center of the ball.
		Here, Beta denotes the angle between the x-axis liine and the balls direction of travel
		The direction of travel can be seen as its velocity vector based at the center of the ball
		Alpha denotes the angle between the x-axis line and the line passing through the center of both balls
		Gamma is the difference between these angles, giving us the angle between the line passing through both balls
		and the velocity-vector of one of the balls.
	*/
	alpha1 := math.Atan2(yDiff, xDiff)
	beta1 := math.Atan2(b1.vy, b1.vx)
	gamma1 := beta1 - alpha1

	alpha2 := math.Atan2(-yDiff, -xDiff)
	beta2 := math.Atan2(b2.vy, b2.vx)
	gamma2 := beta2 - alpha2

	// Compute vector norms before collision by splitting the velocity vectors of the balls
	// One on the line between the centers of the balls
	// One perpendicular to it, this size should not be affected
	b1LineNorm := getVelocity(*b1) * math.Cos(gamma1)
	b1PerpendNorm := getVelocity(*b1) * math.Sin(gamma1)
	b2LineNorm := getVelocity(*b2) * math.Cos(gamma2)
	b2PerpendNorm := getVelocity(*b2) * math.Sin(gamma2)

	// Compute affected vectors after the collision
	// Use formulas given by Wolfram Alpha using both K.E and Momentum equations
	b1LineAfterNorm := ((b1.mass-b2.mass)*b1LineNorm - 2*b2.mass*b2LineNorm) / (b1.mass + b2.mass)

	b2LineAfterNorm := (b1.mass-b2.mass)*b2LineNorm + 2*b1.mass*b1LineNorm/(b1.mass+b2.mass)

	// Now we can compute the velocities post collision for both balls
	// both in the x- and y-direction
	postCollisionx1 := (b1PerpendNorm * math.Sin(alpha1) * -1) + (b1LineAfterNorm * math.Cos(alpha1))
	postCollisiony1 := (b1PerpendNorm * math.Cos(alpha1)) + (b1LineAfterNorm * math.Sin(alpha1))

	postCollisionx2 := (b2PerpendNorm * math.Sin(alpha2) * -1) - (b2LineAfterNorm * math.Cos(alpha2))
	postCollisiony2 := (b2PerpendNorm * math.Cos(alpha2)) - (b2LineAfterNorm * math.Sin(alpha2))

	// Set the new velocities in x- and y-directions
	// Also take an extra step, should help prevent the same collision being handled several times
	b1.vx = postCollisionx1
	b1.vy = postCollisiony1
	b1.x += b1.vx * deltaT
	b1.y += b1.vy * deltaT

	b2.vx = postCollisionx2
	b2.vy = postCollisiony2
	b2.x += b2.vx * deltaT
	b2.y += b2.vy * deltaT

}

// Returns the norm of the velocity vector of the Ball b
func getVelocity(b Ball) float64 {
	return math.Sqrt(math.Pow(b.vx, 2) + math.Pow(b.vy, 2))
}
