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
	"sync"
)

type idsequence struct {
	sync.Mutex
	next int
}

func (i *idsequence) Next() int {
	i.Lock()
	defer i.Unlock()

	x := i.next
	i.next += 1
	return x
}

func (i *idsequence) load(genomes []Genome) {
	for _, g := range genomes {
		if i.next <= g.ID {
			i.next = g.ID + 1
		}
		for _, c := range g.Conns {
			if i.next <= c.Innovation {
				i.next = c.Innovation + 1
			}
		}
	}
}

type nkey struct {
	x, y float64
}

type ckey struct {
	source, target int
}

type marker struct {
	ids   IDSequence
	nodes map[nkey]int
	conns map[ckey]int
	sync.Mutex
}

func (m *marker) MarkConn(c *Connection) {
	m.Lock()
	defer m.Unlock()
	k := ckey{c.Source, c.Target}
	if inno, ok := m.conns[k]; ok {
		c.Innovation = inno
	} else {
		c.Innovation = m.ids.Next()
		m.conns[k] = c.Innovation
	}
}

func (m *marker) MarkNode(n *Node) {
	m.Lock()
	defer m.Unlock()
	k := nkey{n.X, n.Y}
	if inno, ok := m.nodes[k]; ok {
		n.Innovation = inno
	} else {
		n.Innovation = m.ids.Next()
		m.nodes[k] = n.Innovation
	}
}

func (m *marker) Reset() {
	m.Lock()
	defer m.Unlock()
	m.nodes = make(map[nkey]int)
	m.conns = make(map[ckey]int)
}

func (m *marker) load(genomes []Genome) {
	m.Reset()
	for _, g := range genomes {
		for _, n := range g.Nodes {
			k := nkey{n.X, n.Y}
			m.nodes[k] = n.Innovation
		}
		for _, c := range g.Conns {
			k := ckey{c.Source, c.Target}
			m.conns[k] = c.Innovation
		}
	}
}
