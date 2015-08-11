package searcher

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/rqme/neat"
	"github.com/rqme/neat/result"
)

type distrec struct {
	id   int
	dist float64
}

type distrecs []distrec

func (d distrecs) Len() int           { return len(d) }
func (d distrecs) Less(i, j int) bool { return d[i].dist < d[j].dist }
func (d distrecs) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

type BehaviorRecord struct {
	ID       int
	Behavior []float64
}

type BehaviorRecords []BehaviorRecord

func (b BehaviorRecords) Len() int           { return len(b) }
func (b BehaviorRecords) Less(i, j int) bool { return b[i].ID < b[j].ID }
func (b BehaviorRecords) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

type NoveltySettings interface {
	NoveltyEvalArchive() bool         // True = evaluate archive during each search
	NoveltyArchiveThreshold() float64 // Threshold at which to admit phenome into archive
	NumNearestNeighbors() int         // K-nearest neighbors
}
type Novelty struct {
	NoveltySettings

	neat.Searcher
	archive   neat.Phenomes
	behaviors BehaviorRecords
	sync.Mutex
}

func (s *Novelty) SetContext(x neat.Context) error {
	x.State()["novelty-behaviors"] = &s.behaviors
	if cx, ok := s.Searcher.(neat.Contextable); ok {
		return cx.SetContext(x)
	}
	return nil
}

// Searches the phenomes one by one and returns the results
func (s Novelty) Search(phenomes []neat.Phenome) ([]neat.Result, error) {

	// Re-evaluate archive phenomes if necessary
	var bs BehaviorRecords = make([]BehaviorRecord, 0, len(phenomes)+len(s.behaviors)) // behaviors
	if !s.NoveltyEvalArchive() && len(s.behaviors) > 0 {
		bs = append(bs, s.behaviors...)
	}

	// Execute the search using the inner searcher
	var rs []neat.Result
	var err error
	if s.NoveltyEvalArchive() && len(s.archive) > 0 {
		rs, err = s.Searcher.Search(append(phenomes, s.archive...))
	} else {
		rs, err = s.Searcher.Search(phenomes)
	}

	if err != nil {
		err = fmt.Errorf("Error running search: %v", err)
		return nil, err
	}

	// Ensure result includes behavior and can SetNovelty
	for _, r := range rs {
		if br, ok := r.(neat.Behaviorable); ok {
			bs = append(bs, BehaviorRecord{ID: r.ID(), Behavior: br.Behavior()})
		} else {
			err = fmt.Errorf("Result from evalutor did not implement Behaviorable")
			return nil, err
		}
	}
	sort.Sort(bs)

	// Calculate novelty
	wg := new(sync.WaitGroup)
	nrs := make([]neat.Result, len(rs))
	for i, r := range rs {
		wg.Add(1)
		go func(idx int, r neat.Result) {

			// Find this phenome's record
			id := r.ID()
			i := sort.Search(len(bs), func(i int) bool { return bs[i].ID >= id })
			bsi := bs[i]

			// Create a list of all the other records
			o := append(bs[:i], bs[i+1:]...)

			// Create the distance records
			var ds distrecs = make([]distrec, len(o))
			for j := 0; j < len(o); j++ {
				ds[j].id = o[j].ID
				ds[j].dist = calcDist(bsi.Behavior, o[j].Behavior)
			}

			// Sort the records and take the K nearest neighbors
			sort.Sort(ds)
			sum := 0.0
			k := s.NumNearestNeighbors()
			for j := 0; j < k; j++ {
				sum += 1.0 / float64(k) * ds[j].dist
			}

			// Set the novelty
			var ok bool
			var nr *result.Novelty
			if nr, ok = r.(*result.Novelty); !ok {
				nr = result.NewNovelty(id, r.Fitness(), r.Err(), r.Stop(), bsi.Behavior)
			}
			nr.SetNovelty(sum)
			nrs[idx] = nr
			if sum > s.NoveltyArchiveThreshold() {
				s.Lock()
				//s.Archive = append(s.Archive, )
				s.behaviors = append(s.behaviors, bsi)
				s.Unlock()
			}
			wg.Done()
		}(i, r)
	}
	wg.Wait()

	// Determine the novelty of each phenome against the
	// Process the results, setting novelty, replace Fitness with novelty
	return nrs, err
}

func calcDist(a, b []float64) float64 {
	sum := 0.0
	for i := 0; i < len(a); i++ {
		sum += ((a[i] - b[i]) * (a[i] - b[i]))
	}
	return math.Sqrt(sum)
}
