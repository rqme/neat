package neat

type Context interface {
	// Component helpers
	Archiver() Archiver
	Comparer() Comparer
	Crosser() Crosser
	Decoder() Decoder
	Evaluator() Evaluator
	Generator() Generator
	Mutator() Mutator
	Searcher() Searcher
	Speciater() Speciater
	Visualizer() Visualizer

	// State is a registry of elements to be persisted
	State() map[string]interface{}

	// Returns the next ID in the sequence
	NextID() int

	// Returns the innovation number for the gene
	Innovation(t InnoType, k InnoKey) int
}
