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
)

// Definition of a synapse
type Connection struct {
	Innovation     int     // Innovation number for this connection
	Source, Target int     // Innovation numbers of the source and target nodes
	Weight         float64 // Connection weight
	Enabled        bool    // Is this connection enabled?
}

func (c Connection) Key() (k InnoKey) {
	k[0] = float64(c.Source)
	k[1] = float64(c.Target)
	return
}

func (c Connection) String() string {
	b := bytes.NewBufferString(fmt.Sprintf("Conn %d Source %d Target %d Weight %f ", c.Innovation, c.Source, c.Target, c.Weight))
	if c.Enabled {
		b.WriteString("Enabled")
	} else {
		b.WriteString("Disabled")
	}
	return b.String()
}

type Connections map[int]Connection

func (cm Connections) connsToSlice() []Connection {

	// Maps with non-string keys cannot be encoded. Transfer to a slice to handle this
	items := &sortConnsByInnovation{make([]Connection, 0, len(cm))}
	for _, s := range cm {
		items.conns = append(items.conns, s)
	}
	sort.Sort(items)
	return items.conns
}

func (cm *Connections) connsFromSlice(items []Connection) {
	if *cm == nil {
		*cm = make(map[int]Connection)
	}
	m := *cm
	for _, s := range items {
		m[s.Innovation] = s
	}
}

func (cm Connections) MarshalJSON() ([]byte, error) {
	items := cm.connsToSlice()
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(items)
	return buf.Bytes(), err
}

func (cm *Connections) UnmarshalJSON(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)
	var items []Connection
	err := dec.Decode(&items)
	if err == nil {
		cm.connsFromSlice(items)
	}
	return err
}

type sortConnsByKey struct {
	nodeMap map[int]Node
	conns   []Connection
}

func (g *sortConnsByKey) Len() int { return len(g.conns) }
func (g *sortConnsByKey) Less(i, j int) bool {
	ti := g.nodeMap[g.conns[i].Target]
	tj := g.nodeMap[g.conns[j].Target]
	if ti.Y == tj.Y {
		return ti.X < tj.X
	} else {
		return ti.Y < tj.Y
	}
}
func (g *sortConnsByKey) Swap(i, j int) { g.conns[i], g.conns[j] = g.conns[j], g.conns[i] }

type sortConnsByInnovation struct{ conns []Connection }

func (g *sortConnsByInnovation) Len() int { return len(g.conns) }
func (g *sortConnsByInnovation) Less(i, j int) bool {
	return g.conns[i].Innovation < g.conns[j].Innovation
}
func (g *sortConnsByInnovation) Swap(i, j int) { g.conns[i], g.conns[j] = g.conns[j], g.conns[i] }
