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
	. "github.com/rqme/errors"
	"github.com/rqme/neat"

	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
)

type File struct {
	ArchivePath string
	ArchiveName string
}

// Archives the configuration extracted from an item to a file
func (a File) Archive(item neat.Configurable) error {

	errs := new(Errors)
	for _, suffix := range []string{"config", "state"} {

		// Extract the configuration for this tag
		c, err := neat.Extract(item, fmt.Sprintf("neat.%s", suffix))
		if err != nil {
			errs.Add(fmt.Errorf("archiver.File.Archive - Error extracting for %s : %v", suffix, err))
			continue
		}

		// Ensure the directory
		if _, err := os.Stat(a.ArchivePath); os.IsNotExist(err) {
			if err = os.Mkdir(a.ArchivePath, os.ModePerm); err != nil {
				errs.Add(fmt.Errorf("Could not create archive path %s: %v", a.ArchivePath, err))
			}
		}

		// Identify the path
		var p string
		if a.ArchiveName == "" {
			p = path.Join(a.ArchivePath, fmt.Sprintf("%s.json", suffix))
		} else {
			p = path.Join(a.ArchivePath, fmt.Sprintf("%s-%s.json", a.ArchiveName, suffix))
		}

		// Create the file
		f, err := os.Create(p)
		if err != nil {
			errs.Add(fmt.Errorf("archiver.File.Archive - Error creating file for %s : %v", suffix, err))
			continue
		}

		// Write the config to the file
		_, err = f.WriteString(c)
		if err != nil {
			errs.Add(fmt.Errorf("archiver.File.Archive - Error writing to file for %s : %v", suffix, err))
			continue
		}
		f.Close()

	}

	return errs.Err()
}

// Restores an item from the configuration stored in a file
func (a File) Restore(item neat.Configurable) error {

	errs := new(Errors)
	for _, suffix := range []string{"config", "state"} {

		// Identify the path
		var p string
		if a.ArchiveName == "" {
			p = path.Join(a.ArchivePath, fmt.Sprintf("%s.json", suffix))
		} else {
			p = path.Join(a.ArchivePath, fmt.Sprintf("%s-%s.json", a.ArchiveName, suffix))
		}

		// Open the file
		f, err := os.Open(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Nothing to restore
			}
			errs.Add(fmt.Errorf("archiver.File.Restore - Error opening file for %s : %v", suffix, err))
			continue
		}

		// Read the config to the file
		b := bytes.NewBufferString("")
		r := bufio.NewReader(f)
		for {
			s, err := r.ReadBytes('\n')
			if err != nil && err != io.EOF {
				break
			} else {
				b.Write(s)
				b.WriteString("\n")
				if err != nil && err == io.EOF {
					break
				}
			}
		}
		if err != nil && err != io.EOF {
			errs.Add(fmt.Errorf("archiver.File.Restore - Error reading from file for %s : %v", suffix, err))
			continue
		}
		f.Close()

		// Configure the item
		err = item.Configure(b.String())
		if err != nil {
			errs.Add(fmt.Errorf("archiver.File.Restore - Error Configuring for %s : %v", suffix, err))
			continue
		}

	}
	return errs.Err()
}
