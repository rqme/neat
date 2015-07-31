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

package example

import (
	"flag"
	"fmt"

	. "github.com/rqme/errors"
	"github.com/rqme/neat"
	"github.com/rqme/neat/archiver"
	"github.com/rqme/neat/comparer"
	"github.com/rqme/neat/crosser"
	"github.com/rqme/neat/decoder"
	"github.com/rqme/neat/generator"
	"github.com/rqme/neat/mutator"
	"github.com/rqme/neat/searcher"
	"github.com/rqme/neat/speciater"
	"github.com/rqme/neat/visualizer"
)

// Command-line flags. These are helpful but not required as the user could supply an archiver based on a SQL or NoSQL store.
var (
	ConfigPath  = flag.String("config-path", "", "Path to configuration file to override archive.")
	ArchivePath = flag.String("archive-path", "", "Path to archive files (should be a directory).")
	ArchiveName = flag.String("archive-name", "", "Name for archive file (e.g., <name>-pop.json, <name>-config.json")
	WebPath     = flag.String("web-path", "", "Path to web directory.")
)

// Creates a new experiment using the default helpers but can be overridden using the options.
func NewExperiment(eval neat.Evaluator, options ...func(*neat.Experiment) error) (*neat.Experiment, error) {

	e, err := neat.New(
		// Add the archiver and restorer
		func(e *neat.Experiment) error {
			e.Archiver = &archiver.File{
				ArchivePath: *ArchivePath,
				ArchiveName: *ArchiveName,
			}
			if *ConfigPath != "" {
				e.Restorer = &archiver.File{
					ArchivePath: *ConfigPath,
					ArchiveName: *ArchiveName,
				}
			} else {
				e.Restorer = e.Archiver.(neat.Restorer)
			}
			return nil
		},

		// Add the decoder
		func(e *neat.Experiment) error {
			e.Decoder = &decoder.Classic{}
			return nil
		},

		// Add the searcher
		func(e *neat.Experiment) error {
			e.Searcher = &searcher.Concurrent{Evaluator: eval}
			return nil
		},

		// Add the generator
		func(e *neat.Experiment) error {
			s := &speciater.Dynamic{}
			s.Comparer = &comparer.Classic{}
			m := &mutator.Classic{}
			c := &crosser.Classic{}
			e.Generator = &generator.Classic{
				Crosser:   c,
				Mutator:   m,
				Speciater: s,
			}
			return nil
		},

		// Add the visualizer
		func(e *neat.Experiment) error {
			v := &visualizer.Web{
				Decoder: e.Decoder,
				WebPath: *WebPath,
			}
			v.Reset()
			e.Visualizer = v
			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Could not create example experiment: %v", err)
	}

	// Override with options
	errs := new(Errors)
	for i, option := range options {
		if err = option(e); err != nil {
			errs.Add(fmt.Errorf("Error setting example experiment with option %d: %v", i, err))
		}
	}
	return e, errs.Err()
}
