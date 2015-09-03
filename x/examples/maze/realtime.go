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
	"path"

	"github.com/rqme/neat"
	"github.com/rqme/neat/result"
)

type RTEvaluator struct {
	Evaluator
	History map[int]*Environment
}

func (e *RTEvaluator) SetPhenomes(ps neat.Phenomes) error {
	if len(e.History) == 0 {
		e.History = make(map[int]*Environment, len(ps))
	}
	old := e.History
	e.History = make(map[int]*Environment, len(ps))
	for _, p := range ps {
		v, ok := old[p.ID()]
		if !ok {
			cln := e.Environment.clone()
			v := &cln
			v.init()
		}
		e.History[p.ID()] = v
	}
	return nil
}

func (e *RTEvaluator) Takedown() error {

	// Collect the endpoints
	pts := make([]Point, len(e.History))
	for _, env := range e.History {
		pts = append(pts, env.Hero.Location)
	}

	// Create maze image with all endpints (cumulative accross all iterations)
	var t string
	if e.useTrial {
		t = fmt.Sprintf("-%d", e.trialNum)
	}
	env := e.Environment.clone()
	return showMaze(path.Join(*WorkPath, fmt.Sprintf("maze%s-end-points.svg", t)), env, nil, pts)
}

func (e *RTEvaluator) Evaluate(p neat.Phenome) (r neat.Result) {

	// Retrieve the environment
	env := e.History[p.ID()]
	h := &env.Hero

	// Iterate the maze 1 step
	var err error
	paths := make([]Line, 0, *Steps)
	stop := false

	// Note the start of the path
	a := h.Location

	// Update the hero's location
	var outputs []float64
	inputs := generateNeuralInputs(*env)
	if outputs, err = p.Activate(inputs); err == nil {
		interpretOutputs(env, outputs[0], outputs[1])
		update(env)
		b := h.Location
		paths = append(paths, Line{A: a, B: b})

		// Look for solution
		stop = distanceToTarget(env) < 5.0
	}

	f := 300.0 - distanceToTarget(env)          // fitness
	bh := []float64{h.Location.X, h.Location.Y} // behavior
	r = &Result{Classic: result.New(p.ID(), f, err, stop), behavior: bh}

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
