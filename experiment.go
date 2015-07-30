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
	"fmt"

	. "github.com/rqme/errors"
)

// Experiment provides the definition of how to solve the problem using NEAT
type Experiment struct {
	// Configuration
	Iterations     int         `neat:config`
	Traits         Traits      `neat:config`
	FitnessType    FitnessType `neat:config`
	ExperimentName string      `neat:config`

	// State
	Population Population `neat:state`
	cache      map[int]Phenome
	bestID     int
	best       Genome
	Iteration  int
	Stopped    bool

	// Helpers
	Archiver
	Restorer
	Decoder
	Generator
	Searcher
	Visualizer
	ids IDSequence
	mrk Marker
}

// String returns a description of the experiment
func (e Experiment) String() string {
	return fmt.Sprintf("Experiment %s at iteration %d has best genome %d with fitness %f", e.ExperimentName, e.Iteration, e.bestID, e.best.Fitness)
}

// helpers is a convenience method to enable iterating the helpers
func (e *Experiment) helpers() []interface{} {
	return []interface{}{
		e.Archiver,
		e.Restorer,
		e.Decoder,
		e.Generator,
		e.Searcher,
		e.Visualizer,
	}
}

// Creates a new Experiment using the specified options. See neat/x/examples for an
// example Experiment setup
func New(options ...func(e *Experiment) error) (*Experiment, error) {

	// Configure the new experiment using the options
	errs := new(Errors)
	e := new(Experiment)
	for _, option := range options {
		err := option(e)
		if err != nil {
			errs.Add(fmt.Errorf("Could not create new experiment: %v", err))
		}
	}

	// Return the experiment and any errors
	return e, errs.Err()
}

func verify(e *Experiment) error {
	errs := new(Errors)
	if e.Iterations < 1 {
		errs.Add(fmt.Errorf("neat.experiment.Validate - Invalid value for Iterations: %d", e.Iterations))
	}

	// Validate settings and helpers
	helpers := e.helpers()
	for _, helper := range helpers {
		if v, ok := helper.(Validatable); ok {
			err := v.Validate()
			if err != nil {
				errs.Add(err)
			}
		}
	}

	return errs.Err()
}

func Run(e *Experiment) error {

	// Restore the configuration and, if available, state
	if err := e.Restore(); err != nil {
		return fmt.Errorf("Could not restore the experiment's configuration and, if available, state: %v", err)
	}

	// Ensure this is a valid experiment
	if err := verify(e); err != nil {
		return fmt.Errorf("Verification of experiment failed: %v", err)
	}

	// Reset the IDs to account for loaded state
	resetIDs(e)

	// Iterate the experiment
	for e.Iteration = 0; e.Iteration < e.Iterations; e.Iteration++ {

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
			e.Stopped = true
			break
		}
	}

	// Take one last archive and return
	if err := e.Archive(); err != nil {
		return fmt.Errorf("Could not take last archive of experiment: %v", err)
	}
	if err := visualize(e); err != nil {
		return fmt.Errorf("Could not visualize the experiment for the last time: %v", err)
	}
	return nil
}

// Archives the experiment to more permanent storage
func (e *Experiment) Archive() (err error) {
	return e.Archiver.Archive(e)
}

// Restores the experiment from the archived version
func (e *Experiment) Restore() error {
	return e.Restorer.Restore(e)
}

// Configure updates the settings and state of the experiment
func (e *Experiment) Configure(config string) error {
	errs := new(Errors)
	err := Configure(config, e)
	if err != nil {
		errs.Add(fmt.Errorf("Could not configure the experiment: %s", err))
	}
	for _, helper := range e.helpers() {
		if x, ok := helper.(Configurable); ok {
			err = x.Configure(config)
			if err != nil {
				errs.Add(err)
			}
		}
	}
	return errs.Err()
}

// Resets the ID sequence and innovation history using the current population
func resetIDs(e *Experiment) {
	// Load the ID sequence
	ids := new(idsequence)
	ids.load(e.Population.Genomes)
	e.ids = ids

	// Load the innovations
	mrk := new(marker)
	mrk.ids = e.ids
	mrk.Reset()
	mrk.load(e.Population.Genomes)
	e.mrk = mrk

	// Set the ID and innovation with the helpers
	helpers := e.helpers()
	for _, helper := range helpers {
		if v, ok := helper.(Identifies); ok {
			v.SetIDs(e.ids)
		}
		if v, ok := helper.(Marks); ok {
			v.SetMarker(e.mrk)
		}
	}
}

// Advances the experiment to the next generation
func advance(e *Experiment) error {

	curr := e.Population
	next, err := e.Generate(curr)
	if err != nil {
		return err
	}

	if next.Generation > curr.Generation {
		if err = visualize(e); err != nil {
			return err
		}
		if err = e.Archive(); err != nil {
			return err
		}
		if err = updateSettings(e, e.best); err != nil {
			return err
		}
	}

	e.Population = next
	return nil
}

// visualize informs the Visualizer helper of the current population
func visualize(e *Experiment) error {
	if e.Visualizer == nil {
		return nil
	}
	return e.Visualize(e.Population)
}

// Update the settings based on the traits of a genome
func updateSettings(e *Experiment, g Genome) error {
	cnt := 0
	b := bytes.NewBufferString("{")
	for t, trait := range e.Traits {
		if trait.IsSetting {
			if cnt > 0 {
				b.WriteString(",\n")
			}
			b.WriteString(fmt.Sprintf(`"%s": %f`, trait.Name, g.Traits[t]))
			cnt += 1
		}
	}
	b.WriteString("\n}")
	return e.Configure(b.String())
}

// Updates the cache of phenomes
func updateCache(e *Experiment) (err error) {
	var old map[int]Phenome
	if len(e.cache) == 0 {
		old = make(map[int]Phenome, 0)
	} else {
		old = e.cache
	}
	e.cache = make(map[int]Phenome, len(e.Population.Genomes))

	errs := new(Errors)
	pc := make(chan Phenome)
	cnt := 0
	for _, g := range e.Population.Genomes {
		if p, ok := old[g.ID]; ok {
			e.cache[g.ID] = p
		} else {
			cnt += 1
			go func(g Genome) {
				p, err := e.Decode(g)
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
	m := make(map[int]int, len(e.Population.Genomes))
	for i, g := range e.Population.Genomes {
		m[g.ID] = i
	}

	// Perform the search
	phenomes := make([]Phenome, 0, len(e.cache))
	for _, p := range e.cache {
		phenomes = append(phenomes, p)
	}
	if h, ok := e.Searcher.(Phenomable); ok {
		if err = h.SetPhenomes(phenomes); err != nil {
			return
		}
	}
	if h, ok := e.Searcher.(Setupable); ok {
		if err = h.Setup(); err != nil {
			return
		}
	}
	results, err := e.Search(phenomes)
	if err != nil {
		return
	}
	if h, ok := e.Searcher.(Takedownable); ok {
		if err = h.Takedown(); err != nil {
			return
		}
	}

	// Update the fitnesses
	var best Genome
	errs := new(Errors)
	fit := make([]float64, len(e.Population.Genomes))
	// TODO: make this concurrent
	for _, r := range results {
		i := m[r.ID()]
		if err = r.Err(); err != nil {
			errs.Add(fmt.Errorf("Error updating fitness for genome [%d]: %v", r.ID(), r.Err()))
		}
		e.Population.Genomes[i].Fitness = r.Fitness()
		fit[i] = e.Population.Genomes[i].Fitness
		if e.Population.Genomes[i].Fitness > best.Fitness {
			best = e.Population.Genomes[i]
		}
		stop = stop || r.Stop()
	}

	// Update the best genome
	if errs.Err() == nil {
		if e.FitnessType == AbsoluteFitness {
			if best.Fitness > e.best.Fitness {
				e.best = best
			}
		} else {
			e.best = best
		}
	}
	return stop, errs.Err()
}
