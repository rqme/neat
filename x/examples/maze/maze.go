package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
)

type Point struct {
	X, Y float64
}

// Angle determines the angle of the vector defined by [0,0] to this point
func angle(p Point) float64 {
	if p.X == 0.0 {
		if p.Y > 0.0 {
			return 90.0
		} else {
			return 270.0
		}
	}
	ang := math.Atan(p.Y/p.X) / math.Pi * 180.0
	if p.X > 0 {
		return ang
	} else {
		return ang + 180.0
	}
}

// Rotates a point around another point
func rotate(angle float64, a *Point, p Point) {
	rad := angle / 180.0 * math.Pi
	a.X -= p.X
	a.Y -= p.Y

	ox := a.X
	oy := a.Y
	a.X = math.Cos(rad)*ox - math.Sin(rad)*oy
	a.Y = math.Sin(rad)*ox - math.Cos(rad)*oy

	a.X += p.X
	a.Y += p.Y
}

// Distance between two points
func distance(a, b Point) float64 {
	dx := b.X - a.X
	dy := b.Y - a.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type Line struct {
	A, B Point
}

// Midpoint of the line segment
func midpoint(l Line) Point {
	return Point{
		X: (l.A.X + l.B.X) / 2.0,
		Y: (l.A.Y + l.B.Y) / 2.0,
	}
}

// Intersection btween two line segments if it exists
func intersection(l0, l1 Line) (pt Point, found bool) {

	a, b := l0.A, l0.B
	c, d := l1.A, l1.B

	rtop := (a.Y-c.Y)*(d.X-c.X) - (a.X-c.X)*(d.Y-c.Y)
	rbot := (b.X-a.X)*(d.Y-c.Y) - (b.Y-a.Y)*(d.X-c.X)
	stop := (a.Y-c.Y)*(b.X-a.X) - (a.X-c.X)*(b.Y-a.Y)
	sbot := (b.X-a.X)*(d.Y-c.Y) - (b.Y-a.Y)*(d.X-c.X)

	if rbot == 0 || sbot == 0 {
		// Lines are parallel
		return
	}

	r := rtop / rbot
	s := stop / sbot

	if r > 0 && r < 1 && s > 0 && s < 1 {
		pt.X = a.X + r*(b.X-a.X)
		pt.Y = a.Y + r*(b.Y-a.Y)
		found = true
	}
	return
}

// Distance between line segment and point
func linedist(l Line, n Point) float64 {
	utop := (n.X-l.A.X)*(l.B.X-l.A.X) + (n.Y-l.A.Y)*(l.B.Y-l.A.Y)
	ubot := distance(l.A, l.B)
	ubot *= ubot
	if ubot == 0.0 {
		return 0.0
	}
	u := utop / ubot

	if u < 0 || u > 1 {
		d1 := distance(l.A, n)
		d2 := distance(l.B, n)
		if d1 < d2 {
			return d1
		} else {
			return d2
		}
	}

	p := Point{
		X: l.A.X + u*(l.B.X-l.A.X),
		Y: l.A.Y + u*(l.B.Y-l.A.Y),
	}
	return distance(p, n)
}

// Line segment length
func (l Line) Length() float64 {
	return distance(l.A, l.B)
}

type Character struct {
	RangeFinderAngles []float64
	RadarAngles1      []float64
	RadarAngles2      []float64

	RangeFinders []float64
	Radar        []float64

	Location         Point
	Heading          float64
	Speed            float64
	Radius           float64
	RangeFinderRange float64
	AngularVelocity  float64
}

type Environment struct {
	Hero      Character
	Lines     []Line
	End       Point
	ReachGoal bool
}

func (e Environment) clone() Environment {
	e2 := Environment{
		Hero: Character{
			Heading:           e.Hero.Heading,
			Speed:             e.Hero.Speed,
			Location:          e.Hero.Location,
			RangeFinderAngles: make([]float64, len(e.Hero.RangeFinderAngles)),
			RadarAngles1:      make([]float64, len(e.Hero.RadarAngles1)),
			RadarAngles2:      make([]float64, len(e.Hero.RadarAngles2)),
			RangeFinders:      make([]float64, len(e.Hero.RangeFinders)),
			Radar:             make([]float64, len(e.Hero.Radar)),
			Radius:            e.Hero.Radius,
			RangeFinderRange:  e.Hero.RangeFinderRange,
			AngularVelocity:   e.Hero.AngularVelocity,
		},
		Lines: e.Lines,
		End:   e.End,
	}
	copy(e2.Hero.RangeFinderAngles, e.Hero.RangeFinderAngles)
	copy(e2.Hero.RadarAngles1, e.Hero.RadarAngles1)
	copy(e2.Hero.RadarAngles2, e.Hero.RadarAngles2)
	copy(e2.Hero.RangeFinders, e.Hero.RangeFinders)
	copy(e2.Hero.Radar, e.Hero.Radar)
	return e2
}

func copyEnvironment(e *Environment) *Environment {
	cln := (*e).clone()
	return &cln
}

// Initializes the environment
func (e *Environment) init() {

	// Set up the hero
	e.Hero.RangeFinderAngles = []float64{-90, -45, 0, 45, 90, -180}
	e.Hero.RadarAngles1 = []float64{315, 45, 135, 225}
	e.Hero.RadarAngles2 = []float64{405, 135, 225, 315}
	e.Hero.Radius = 8.0
	e.Hero.RangeFinderRange = 100.0
	e.Hero.RangeFinders = make([]float64, len(e.Hero.RangeFinderAngles))
	e.Hero.Radar = make([]float64, len(e.Hero.RadarAngles1))
}

// String displays debug information
func (e Environment) String() string {
	b := bytes.NewBufferString("")
	b.WriteString(fmt.Sprintf("Hero: %f %f\n", e.Hero.Location.X, e.Hero.Location.Y))
	b.WriteString(fmt.Sprintf("EndPoint: %f %f\n", e.End.X, e.End.Y))
	b.WriteString("Lines\n")
	for _, l := range e.Lines {
		b.WriteString(fmt.Sprintf("\t%f %f %f %f\n", l.A.X, l.A.Y, l.B.X, l.B.Y))
	}
	return b.String()
}

// Used for fitness calculations
func distanceToTarget(e *Environment) float64 {
	dist := distance(e.Hero.Location, e.End)
	if math.IsNaN(dist) {
		log.Println("NAN distance error...") // Should this be an actual error?
		return 500.0
	}
	if dist < 5.0 {
		e.ReachGoal = true
	}
	return dist
}

// Create neural net inputs from sensors
func generateNeuralInputs(e Environment) (inputs []float64) {
	inputs = make([]float64, 10)

	// range finders
	n := len(e.Hero.RangeFinders)
	for i := 0; i < n; i++ {
		inputs[i] = e.Hero.RangeFinders[i] / e.Hero.RangeFinderRange
		if math.IsNaN(inputs[i]) {
			log.Println("NAN in inputs")
		}
	}

	// Radar
	for i := 0; i < len(e.Hero.Radar); i++ {
		inputs[i+n] = e.Hero.Radar[i]
		if math.IsNaN(inputs[i+n]) {
			log.Println("NAN in inputs")
		}
	}
	return
}

// Transforms neural net outputs into angular velocity and speed
func interpretOutputs(e *Environment, o1, o2 float64) {
	if math.IsNaN(o1) || math.IsNaN(o2) {
		log.Println("NAN in outputs")
	}

	e.Hero.AngularVelocity += (o1 - 0.5) * 1.0
	e.Hero.Speed += (o2 - 0.5) * 1.0

	// constraints of speed and angular velocity
	if e.Hero.Speed > 3.0 {
		e.Hero.Speed = 3.0
	} else if e.Hero.Speed < -3.0 {
		e.Hero.Speed = -3.0
	}
	if e.Hero.AngularVelocity > 3.0 {
		e.Hero.AngularVelocity = 3.0
	} else if e.Hero.AngularVelocity < -3.0 {
		e.Hero.AngularVelocity = -3.0
	}
}

// Run a time step of the simulation
func update(e *Environment) {
	if e.ReachGoal {
		return
	}
	vx := math.Cos(e.Hero.Heading/180.0*math.Pi) * e.Hero.Speed
	vy := math.Sin(e.Hero.Heading/180.0*math.Pi) * e.Hero.Speed
	if math.IsNaN(vx) {
		log.Println("vx is NAN")
	}

	e.Hero.Heading += e.Hero.AngularVelocity
	if math.IsNaN(e.Hero.AngularVelocity) {
		log.Println("Hero AngularVelocity is NAN")
	}

	if e.Hero.Heading > 360 {
		e.Hero.Heading -= 360
	} else if e.Hero.Heading < 0 {
		e.Hero.Heading += 360
	}

	newloc := Point{
		X: vx + e.Hero.Location.X,
		Y: vy + e.Hero.Location.Y,
	}

	// collision detection
	if !collideLines(e, newloc) {
		e.Hero.Location = newloc
	}

	updateRangeFinders(e)
	updateRadar(e)
}

// See if navigator has hit anything
func collideLines(e *Environment, loc Point) bool {
	for _, l := range e.Lines {
		if linedist(l, loc) < e.Hero.Radius {
			return true
		}
	}
	return false
}

// Range finder sensors
func updateRangeFinders(e *Environment) {
	h := &e.Hero
	for i := 0; i < len(h.RangeFinders); i++ {
		rad := h.RangeFinderAngles[i] / 180.0 * math.Pi // radians

		// Project a point from the hero's location outwards
		projPoint := Point{
			X: h.Location.X + math.Cos(rad)*h.RangeFinderRange,
			Y: h.Location.Y + math.Sin(rad)*h.RangeFinderRange,
		}

		// Rotate the projected point by the hero's heading
		rotate(h.Heading, &projPoint, h.Location)

		// Create a line segment from the hero's location to projected
		projectedLine := Line{A: h.Location, B: projPoint}

		rnge := h.RangeFinderRange // Set range to max by default

		// Now test against the environment to see if we hit anything
		for _, line := range e.Lines {
			if intpt, found := intersection(line, projectedLine); found {
				// If so, then update the range to the distance
				foundRange := distance(intpt, h.Location)
				if foundRange < rnge {
					rnge = foundRange
				}
			}
		}
		if math.IsNaN(rnge) {
			log.Println("Range is NAN")
		}
		h.RangeFinders[i] = rnge
	}
}

func updateRadar(e *Environment) {
	h := &e.Hero
	target := e.End

	// Rotate goal with respect to heading of navigator
	rotate(-h.Heading, &target, h.Location)

	// Translate with respect to location of navigator
	target.X -= h.Location.X
	target.Y -= h.Location.Y

	// What angle is the vector between target and navigator
	a := angle(target)

	// Fire the appropriate radar sensor
	for i := 0; i < len(h.RadarAngles1); i++ {
		h.Radar[i] = 0.0
		if a >= h.RadarAngles1[i] && a <= h.RadarAngles2[i] {
			h.Radar[i] = 1.0
		} else if a+360.0 >= h.RadarAngles1[i] && a+360. <= h.RadarAngles2[i] {
			h.Radar[i] = 1.0
		}
	}
}
