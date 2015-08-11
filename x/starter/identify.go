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

func newIdentify() identify {
	return identify{
		innos: make(map[innovation]int, 100),
	}
}

// NextID returns the next id in the context's sequence
func (x *identify) NextID() int {
	x.Lock()
	defer x.Unlock()
	x.lastID += 1
	return x.lastID
}

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
