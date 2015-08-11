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

package archiver

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/rqme/neat"
)

type FileSettings interface {
	ArchiveName() string
	ArchivePath() string
}

type File struct {
	FileSettings
	useTrials bool
	trialNum  int
}

func (a *File) SetTrial(t int) error {
	a.useTrials = true
	a.trialNum = t
	return nil
}

func (a *File) makePath(s string) string {
	p := a.ArchivePath()
	if a.useTrials {
		p = path.Join(p, strconv.Itoa(a.trialNum))
	}
	return path.Join(p, fmt.Sprintf("%s-%s.json", a.ArchiveName(), s))
}

func (a *File) Archive(ctx neat.Context) error {

	// Save the settings
	name := a.makePath("config")
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	e := json.NewEncoder(f)
	if err = e.Encode(ctx); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// Save the state values
	for k, v := range ctx.State() {
		name := a.makePath(k)
		f, err = os.Create(name)
		if err != nil {
			return err
		}
		e = json.NewEncoder(f)
		if err = e.Encode(v); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

func (a *File) Restore(ctx neat.Context) error {

	// Restore the settings
	name := a.makePath("config")
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	d := json.NewDecoder(f)
	if err = d.Decode(&ctx); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// Restore the state values
	for k, v := range ctx.State() {
		name := a.makePath(k)
		if _, err := os.Stat(name); os.IsNotExist(err) {
			continue
		}

		f, err = os.Open(name)
		if err != nil {
			return err
		}
		d = json.NewDecoder(f)
		if err = d.Decode(&v); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}
