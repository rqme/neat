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
	. "github.com/rqme/errors"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

// Configures a item based on the JSON-encoded string
// TODO: Restrict changes to only fields that are tagged
func Configure(cfg string, item interface{}) error {
	b := bytes.NewBufferString(cfg)
	d := json.NewDecoder(b)
	err := d.Decode(&item)
	if err != io.EOF {
		return err
	}
	return nil
}

// Extract exctacts the fields with a specific tag into a JSON string
func Extract(item interface{}, section string) (string, error) {

	// Extract the configuration
	errs := new(Errors)
	config := make(map[string]string, 100)
	err := extractFields(config, item, section)
	if err != nil {
		errs.Add(err)
	}

	// Alphabetize the keys
	ks := make([]string, 0, len(config))
	for k, _ := range config {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	// Create the string
	b := bytes.NewBufferString("")
	i := 0
	if len(config) > 0 {
		b.WriteString(`{`)
		for _, k := range ks {
			v := config[k]
			b.WriteString(fmt.Sprintf(`"%s": %s`, k, v))
			if i < len(config)-1 {
				b.WriteString(", ")
			}
			i += 1
		}
		b.WriteString(`}`)
	}
	return b.String(), errs.Err()
}

func extractFields(cfg map[string]string, item interface{}, section string) error {
	errs := new(Errors)
	v := reflect.ValueOf(item)
	switch v.Kind() {
	case reflect.Struct:

		t := v.Type()
		// Iterate fields. If a field is a Configurable instance, remove its entry and merge in
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).PkgPath != "" {
				continue // This is a private member
			}
			if t.Field(i).Tag.Get("neat") == section {
				if v.Field(i).Interface() != nil {
					b := bytes.NewBufferString("")
					j := json.NewEncoder(b)
					if err := j.Encode(v.Field(i).Interface()); err != nil {
						//DBG("configure.extractFields: name %s value %v error %v", t.Field(i).Name, v.Field(i).Interface(), err)
						errs.Add(err)
					} else {
						s := strings.TrimSpace(b.String())
						if s != "null" && s != "[]" {
							cfg[t.Field(i).Name] = s
						}
					}
				}
			} else {
				if v.Field(i).Kind() == reflect.Struct || v.Field(i).Kind() == reflect.Ptr || v.Field(i).Kind() == reflect.Interface {
					e2 := extractFields(cfg, v.Field(i).Interface(), section)
					if e2 != nil {
						errs.Add(e2)
					}
				}
			}
		}

	case reflect.Ptr:
		return extractFields(cfg, v.Elem().Interface(), section)

	default:
		//me.AddMessage("neat.configure.extractfields - Unsupported kind: %v", v.Kind())
	}

	return errs.Err()
}
