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

package searcher

import (
	"github.com/rqme/neat"
)

type Concurrent struct {
	neat.Evaluator
}

// Configures the helper from a JSON string
func (s *Concurrent) Configure(cfg string) error {
	if x, ok := s.Evaluator.(neat.Configurable); ok {
		return x.Configure(cfg)
	}
	return nil
}

// Searches the phenomes concurrently and returns the results
func (s Concurrent) Search(phenomes []neat.Phenome) ([]neat.Result, error) {
	r := make(chan neat.Result)
	for _, p := range phenomes {
		go func(p neat.Phenome) {
			r <- s.Evaluate(p)
		}(p)
	}
	results := make([]neat.Result, len(phenomes))
	for i := 0; i < len(phenomes); i++ {
		results[i] = <-r
	}
	return results, nil
}

func (s *Concurrent) Setup() error {
	if h, ok := s.Evaluator.(neat.Setupable); ok {
		return h.Setup()
	}
	return nil
}

func (s *Concurrent) Takedown() error {
	if h, ok := s.Evaluator.(neat.Takedownable); ok {
		return h.Takedown()
	}
	return nil
}

func (s *Concurrent) SetPhenomes(p neat.Phenomes) error {
	if h, ok := s.Evaluator.(neat.Phenomable); ok {
		return h.SetPhenomes(p)
	}
	return nil
}
