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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLock(t *testing.T) {
	td, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(td)

	path := filepath.Join(td, "foo.lock")

	lock1 := NewFlock(path)
	lock2 := NewFlock(path)

	if e := lock1.Lock(); e != nil || !lock1.Locked() {
		t.Errorf("lock1 times 1 error:%v,%v", lock1.Locked(), e)
	}
	if e := lock1.Lock(); e != nil || !lock1.Locked() {
		t.Errorf("lock1 times 2 error:%v,%v", lock1.Locked(), e)
	}
	if e := lock1.TryLock(); e != nil || !lock1.Locked() {
		t.Errorf("trylock1 times 1 error:%v,%v", lock1.Locked(), e)
	}
	if e := lock1.TryLock(); e != nil || !lock1.Locked() {
		t.Errorf("trylock1 times 2 error:%v,%v", lock1.Locked(), e)
	}

	if e := lock2.Lock(); e == nil || lock2.Locked() {
		t.Errorf("lock2 times 1 error:%v,%v", lock2.Locked(), e)
	}
	if e := lock2.Lock(); e == nil || lock2.Locked() {
		t.Errorf("lock2 times 2 error:%v,%v", lock2.Locked(), e)
	}
	if e := lock2.TryLock(); e == nil || lock2.Locked() {
		t.Errorf("trylock2 times 1 error:%v,%v", lock2.Locked(), e)
	}
	if e := lock2.TryLock(); e == nil || lock2.Locked() {
		t.Errorf("trylock2 times 2 error:%v,%v", lock2.Locked(), e)
	}

	if e := lock1.Unlock(); e != nil || lock1.Locked() {
		t.Errorf("unlock1 times 1 error:%v,%v", lock1.Locked(), e)
	}

	if e := lock1.Unlock(); e == nil || lock1.Locked() {
		t.Errorf("unlock1 times 2 error:%v,%v", lock1.Locked(), e)
	}

	if e := lock2.Lock(); e != nil || !lock2.Locked() {
		t.Errorf("lock2 times 1 error:%v,%v", lock2.Locked(), e)
	}

	if e := lock2.Lock(); e != nil || !lock2.Locked() {
		t.Errorf("lock2 times 1 error:%v,%v", lock2.Locked(), e)
	}

	if e := lock2.Unlock(); e != nil || lock2.Locked() {
		t.Errorf("unlock2 times 1 error:%v,%v", lock2.Locked(), e)
	}
}
