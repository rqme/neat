/*
Copyright (c) 2015, Brian Hummer (brian@redq.me)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"fmt"
	"log"
	"os"
	"path"

	svg "github.com/ajstarks/svgo"
	"github.com/rqme/neat"
	"github.com/rqme/neat/result"
)

type Evaluator struct {
	Environment
	show bool

	useTrial bool
	trialNum int
	points   []Point
}

func (e *Evaluator) ShowWork(s bool) {
	e.show = s
}

func (e *Evaluator) SetTrial(t int) error {
	e.useTrial = true
	e.trialNum = t
	return nil
}

func (e *Evaluator) Takedown() error {
	// Create maze image with all endpints (cumulative accross all iterations)
	var t string
	if e.useTrial {
		t = fmt.Sprintf("-%d", e.trialNum)
	}
	env := e.Environment.clone()
	return showMaze(path.Join(*WorkPath, fmt.Sprintf("maze%s-end-points.svg", t)), env, nil, e.points)
}

func (e *Evaluator) Evaluate(p neat.Phenome) (r neat.Result) {

	// Clone the environemnt
	var env *Environment
	cln := e.Environment.clone()
	env = &cln
	env.init()
	h := &env.Hero

	// Iterate the maze
	var err error
	paths := make([]Line, 0, *Steps)
	stop := false
	for i := 0; i < *Steps; i++ {

		// Note the start of the path
		a := h.Location

		// Update the hero's location
		var outputs []float64
		inputs := generateNeuralInputs(*env)
		if outputs, err = p.Activate(inputs); err != nil {
			break
		}
		interpretOutputs(env, outputs[0], outputs[1])
		update(env)
		b := h.Location
		paths = append(paths, Line{A: a, B: b})

		// Look for solution
		stop = distanceToTarget(env) < 5.0
		if stop {
			break
		}
	}
	f := 300.0 - distanceToTarget(env) // fitness
	e.points = append(e.points, h.Location)
	b := []float64{h.Location.X, h.Location.Y} // behavior
	r = &Result{Classic: result.New(p.ID(), f, err, stop), behavior: b}

	// Output the maze
	if e.show {
		// Write the file
		var t string
		if e.useTrial {
			t = fmt.Sprintf("-%d", e.trialNum)
		}
		if err := showMaze(path.Join(*WorkPath, fmt.Sprintf("maze%s-%d.svg", t, p.ID())), *env, paths, nil); err != nil {
			log.Println("Could not output maze run:", err)
		}
	}
	return
}

func showMaze(p string, e Environment, hist []Line, pts []Point) error {

	// Determine image size
	var h, w float64
	for _, line := range e.Lines {
		if line.A.X > w {
			w = line.A.X
		}
		if line.A.Y > h {
			h = line.A.Y
		}
		if line.B.X > w {
			w = line.B.X
		}
		if line.B.Y > h {
			h = line.B.Y
		}
	}

	// Create the image
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	img := svg.New(f)
	img.Start(int(w), int(h))
	defer img.End()

	// Add the maze
	if len(hist) > 0 {
		img.Circle(int(hist[0].A.X), int(hist[0].A.Y), 4, `fill="green"`) // start
	}
	img.Circle(int(e.End.X), int(e.End.Y), 4, `fill="red"`)

	for _, line := range e.Lines {
		img.Path(fmt.Sprintf("M %f %f L %f %f", line.A.X, line.A.Y, line.B.X, line.B.Y), `stroke-width="1" stroke="black" fill="none"`)
	}

	for _, line := range hist {
		img.Path(fmt.Sprintf("M %f %f L %f %f", line.A.X, line.A.Y, line.B.X, line.B.Y), `stroke-width="1" stroke="blue" fill="none"`)

	}

	for _, point := range pts {
		img.Circle(int(point.X), int(point.Y), 1, `fill="green"`)
	}
	return nil
}
