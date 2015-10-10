/*
Copyright 2013 The Go Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lock

import (
	"os"
	"path/filepath"
	"sync"
)

// Flock Is not test on plan9
type flock struct {
	path    string
	absPath string
	mu      sync.RWMutex
	fh      *os.File
	locked  bool
}

func NewFlock(path string) Flocker {
	f := &flock{path: path}
	f.absPath, _ = filepath.Abs(path)
	return f
}

func (f *flock) Path() string {
	return f.path
}

func (f *flock) Locked() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.locked
}

func (f *flock) String() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.locked {
		return fmt.Sprintf("'%s' locked", f.path)
	} else {
		return fmt.Sprintf("'%s' unlock", f.path)
	}
}

func (f *flock) Lock() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.locked {
		return nil
	}

	fh, err := os.OpenFile(f.absPath, os.O_RDWR|os.O_CREATE, os.ModeExclusive|0644)
	if err != nil {
		return err
	}

	f.locked = err == nil
	if f.locked {
		f.fh = fh
	}
	return err
}

func (f *flock) Unlock() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.locked {
		return ErrUnlock
	}

	err := f.fh.Close()
	if err == nil {
		f.locked = false
		f.fh = nil
	}
	return err
}
