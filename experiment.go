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

package neat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	. "github.com/rqme/errors"
)

type ExperimentSettings interface {
	Iterations() int
	Traits() Traits
	FitnessType() FitnessType
	ExperimentName() string
}

// Experiment provides the definition of how to solve the problem using NEAT
type Experiment struct {
	ExperimentSettings
	ctx Context

	// State
	population Population `neat:"state"`
	cache      map[int]Phenome
	best       Genome
	iteration  int
	stopped    bool
}

func (e *Experiment) SetContext(x Context) error {
	e.ctx = x
	e.ctx.State()["population"] = &e.population
	return nil
}

func (e Experiment) Context() Context { return e.ctx }

func (e Experiment) Population() Population { return e.population }

func (e Experiment) Stopped() bool { return e.stopped }

func (e Experiment) Iteration() int { return e.iteration }

// String returns a description of the experiment
func (e Experiment) String() string {
	return fmt.Sprintf("Experiment %s at iteration %d has best genome %d with fitness %f", e.ExperimentName(), e.iteration, e.best.ID, e.best.Fitness)
}

// Runs a configured experiment. If restoring, including just the configuration, this must be done
// prior to calling Run.
func Run(e *Experiment) error {

	// Ensure this is a valid experiment
	if e.Iterations() < 1 {
		return fmt.Errorf("Invalid value for Iterations: %d", e.Iterations())
	}

	// Iterate the experiment
	for e.iteration = 0; e.iteration < e.Iterations(); e.iteration++ {
		//fmt.Println("iteration", e.iteration, "best", e.best.Fitness)
		// Reset the innovation history
		//e.mrk.Reset()

		// Advance the population
		if err := advance(e); err != nil {
			return fmt.Errorf("Could not advance the population: %v", err)
		}

		// Update the phenome cache
		if err := updateCache(e); err != nil {
			return fmt.Errorf("Couuld not update cache in the experiment: %v", err)
		}

		// Evaluate the population
		if stop, err := search(e); err != nil {
			return fmt.Errorf("Error evaluating the population: %v", err)
		} else if stop {
			e.stopped = true
			break
		}
	}

	// Take one last archive and return
	if err := e.ctx.Archiver().Archive(e.ctx); err != nil {
		return fmt.Errorf("Could not take last archive of experiment: %v", err)
	}
	if err := e.ctx.Visualizer().Visualize(e.population); err != nil {
		return fmt.Errorf("Could not visualize the experiment for the last time: %v", err)
	}
	return nil
}

// Advances the experiment to the next generation
func advance(e *Experiment) error {

	curr := e.population
	next, err := e.ctx.Generator().Generate(curr)
	if err != nil {
		return err
	}

	if next.Generation > curr.Generation {
		if err = e.ctx.Visualizer().Visualize(e.population); err != nil {
			return err
		}
		if err = e.ctx.Archiver().Archive(e.ctx); err != nil {
			return err
		}
		if err = updateSettings(e, e.best); err != nil {
			return err
		}
	}

	e.population = next
	return nil
}

// Update the settings based on the traits of a genome
func updateSettings(e *Experiment, g Genome) error {
	cnt := 0
	b := bytes.NewBufferString("{")
	for t, trait := range e.Traits() {
		if trait.IsSetting {
			if cnt > 0 {
				b.WriteString(",\n")
			}
			b.WriteString(fmt.Sprintf(`"%s": %f`, trait.Name, g.Traits[t]))
			cnt += 1
		}
	}
	b.WriteString("\n}")

	enc := json.NewEncoder(b)
	return enc.Encode(&e.ctx)
}

// Updates the cache of phenomes
func updateCache(e *Experiment) (err error) {
	var old map[int]Phenome
	if len(e.cache) == 0 {
		old = make(map[int]Phenome, 0)
	} else {
		old = e.cache
	}
	e.cache = make(map[int]Phenome, len(e.population.Genomes))

	errs := new(Errors)
	pc := make(chan Phenome)
	cnt := 0
	for _, g := range e.population.Genomes {
		if p, ok := old[g.ID]; ok {
			e.cache[g.ID] = p
		} else {
			cnt += 1
			go func(g Genome) {
				p, err := e.ctx.Decoder().Decode(g)
				if err != nil {
					errs.Add(fmt.Errorf("Unable to decode genome [%d]: %v", g.ID, err))
				}
				pc <- p
			}(g)
		}
	}

	for i := 0; i < cnt; i++ {
		p := <-pc
		if p != nil {
			e.cache[p.ID()] = p
		}

	}
	return errs.Err()
}

// Searches the population and updates the genomes' fitness
func search(e *Experiment) (stop bool, err error) {

	// Map the genomes for convenience
	m := make(map[int]int, len(e.population.Genomes))
	for i, g := range e.population.Genomes {
		m[g.ID] = i
	}

	// Perform the search
	phenomes := make([]Phenome, 0, len(e.cache))
	for _, p := range e.cache {
		phenomes = append(phenomes, p)
	}
	for _, h := range []interface{}{e.ctx.Searcher(), e.ctx.Evaluator()} {
		if ph, ok := h.(Phenomable); ok {
			if err = ph.SetPhenomes(phenomes); err != nil {
				return
			}
		}
		if sh, ok := h.(Setupable); ok {
			if err = sh.Setup(); err != nil {
				return
			}
		}
	}

	var rs Results
	if rs, err = e.ctx.Searcher().Search(phenomes); err != nil {
		return
	}

	for _, h := range []interface{}{e.ctx.Evaluator(), e.ctx.Searcher()} {
		if th, ok := h.(Takedownable); ok {
			if err = th.Takedown(); err != nil {
				return
			}
		}
	}

	// Update the fitnesses
	var best Genome
	errs := new(Errors)
	// := make([]float64, len(e.population.Genomes))
	// TODO: make this concurrent
	for _, r := range rs {
		i := m[r.ID()]
		if err = r.Err(); err != nil {
			errs.Add(fmt.Errorf("Error updating fitness for genome [%d]: %v", r.ID(), r.Err()))
		}
		e.population.Genomes[i].Fitness = r.Fitness()
		if imp, ok := r.(Improvable); ok {
			e.population.Genomes[i].Improvement = imp.Improvement()
		} else {
			e.population.Genomes[i].Improvement = e.population.Genomes[i].Fitness
		}
		//fit[i] = e.population.Genomes[i].Fitness
		if e.population.Genomes[i].Fitness > best.Fitness {
			best = e.population.Genomes[i]
		}
		stop = stop || r.Stop()
	}

	// Update the best genome
	if errs.Err() == nil {
		if e.FitnessType() == Absolute {
			if best.Fitness > e.best.Fitness {
				e.best = best
			}
		} else {
			e.best = best
		}
	}

	// Leave the genomes sorted by their fitness descending
	sort.Sort(sort.Reverse(e.population.Genomes))
	return stop, errs.Err()
}
