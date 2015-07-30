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

package mutator

import (
	. "github.com/rqme/errors"
	"github.com/rqme/neat"
)

type Classic struct {
	Complexify
	Weight
	Trait
}

// Configures the helper from a JSON string
func (m *Classic) Configure(cfg string) error {
	errs := new(Errors)
	err := m.Complexify.Configure(cfg)
	if err != nil {
		errs.Add(err)
	}
	err = m.Weight.Configure(cfg)
	if err != nil {
		errs.Add(err)
	}
	err = m.Trait.Configure(cfg)
	if err != nil {
		errs.Add(err)
	}
	return errs.Err()
}

func (c Classic) Mutate(g *neat.Genome) error {
	errs := new(Errors)
	old := g.Complexity()
	err := c.Complexify.Mutate(g)
	if err != nil {
		errs.Add(err)
	}
	if g.Complexity() == old {
		if err = c.Weight.Mutate(g); err != nil {
			errs.Add(err)
		}
		if err = c.Trait.Mutate(g); err != nil {
			errs.Add(err)
		}
	}
	return errs.Err()
}

// Sets the marker for recording innovations
func (m *Classic) SetMarker(marker neat.Marker) { m.Complexify.SetMarker(marker) }
