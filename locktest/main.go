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

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/tgulacsi/lock"
)

func main() {
	t := new(testing.T)
	f := os.Getenv("TEST_LOCK_FILE")
	//log.Printf("f=%s", f)
	if f == "" {
		testLock(t)
	} else {
		testLockInChild(t)
	}
}

func testLockInChild(t *testing.T) {
	f := os.Getenv("TEST_LOCK_FILE")
	log.Printf("testLockInChild f=%s", f)
	if f == "" {
		// not child
		return
	}

	log.Printf("locking %s", f)
	lk, err := lock.Lock(f)
	if err != nil {
		log.Fatalf("Lock failed: %v", err)
	}

	if v, _ := strconv.ParseBool(os.Getenv("TEST_LOCK_CRASH")); v {
		// Simulate a crash, or at least not unlocking the
		// lock.  We still exit 0 just to simplify the parent
		// process exec code.
		os.Exit(0)
	}
	lk.Close()
}

func testLock(t *testing.T) {
	td, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(td)

	path := filepath.Join(td, "foo.lock")

	childLock := func(crash bool) error {
		//log.Printf("executing %s", os.Args[0])
		cmd := exec.Command(os.Args[0])
		cmd.Env = []string{"TEST_LOCK_FILE=" + path}
		if crash {
			cmd.Env = append(cmd.Env, "TEST_LOCK_CRASH=1")
		}
		out, err := cmd.CombinedOutput()
		log.Printf("Child output: %q (err %v)", out, err)
		if err != nil {
			return fmt.Errorf("Child Process lock of %s failed: %v %s", path, err, out)
		}
		return nil
	}

	log.Printf("Locking in crashing child...")
	if err := childLock(true); err != nil {
		log.Fatalf("first lock in child process: %v", err)
	}

	log.Printf("Locking+unlocking in child...")
	if err := childLock(false); err != nil {
		log.Fatalf("lock in child process after crashing child: %v", err)
	}

	log.Printf("Locking in parent...")
	lk1, err := lock.Lock(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Again in parent...")
	_, err = lock.Lock(path)
	if err == nil {
		log.Fatal("expected second lock to fail")
	}

	log.Printf("Locking in child...")
	if childLock(false) == nil {
		log.Fatalf("expected lock in child process to fail")
	}

	log.Printf("Unlocking lock in parent")
	if err := lk1.Close(); err != nil {
		log.Fatal(err)
	}

	lk3, err := lock.Lock(path)
	if err != nil {
		log.Fatal(err)
	}
	lk3.Close()
}
