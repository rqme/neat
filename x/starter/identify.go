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
package starter

import (
	"sync"

	"github.com/rqme/neat"
)

type innovation struct {
	Type neat.InnoType
	Key  neat.InnoKey
}

type identify struct {
	sync.Mutex
	lastID int
	innos  map[innovation]int
}

// NextID returns the next id in the context's sequence
func (x *identify) NextID() int {
	x.Lock()
	defer x.Unlock()
	x.lastID += 1
	return x.lastID
}

// Returns the innovation number for this type and key
func (x *identify) Innovation(t neat.InnoType, k neat.InnoKey) int {
	x.Lock()
	defer x.Unlock()
	var id int
	var ok bool

	in := innovation{Type: t, Key: k}
	if id, ok = x.innos[in]; !ok {
		x.lastID += 1
		id = x.lastID
		x.innos[in] = id
	}
	return id
}

// Initializes the ID sequence and innovation history
func (x *identify) SetPopulation(p neat.Population) error {
	for _, g := range p.Genomes {
		if g.ID > x.lastID {
			x.lastID = g.ID
		}
		for _, n := range g.Nodes {
			if n.Innovation > x.lastID {
				x.lastID = n.Innovation
			}
			x.innos[innovation{Type: neat.NodeInnovation, Key: n.Key()}] = n.Innovation
		}
		for _, c := range g.Conns {
			if c.Innovation > x.lastID {
				x.lastID = c.Innovation
			}
			x.innos[innovation{Type: neat.ConnInnovation, Key: c.Key()}] = c.Innovation
		}
	}
	return nil
}
