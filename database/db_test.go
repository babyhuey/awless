/*
Copyright 2017 WALLIX

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

package database

import (
	"os"
	"testing"
	"time"
)

func TestGetSetDatabaseValues(t *testing.T) {
	db, close := newTestDB()
	defer close()

	value, e := db.GetStringValue("mykey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := value, ""; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	e = db.SetStringValue("mykey", "myvalue")
	if e != nil {
		t.Fatal(e)
	}

	value, e = db.GetStringValue("mykey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := value, "myvalue"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	i, e := db.GetIntValue("myintkey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := i, 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	e = db.SetIntValue("myintkey", 10)
	if e != nil {
		t.Fatal(e)
	}

	i, e = db.GetIntValue("myintkey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := i, 10; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	stamp, e := db.GetTimeValue("mytimekey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := stamp.IsZero(), true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}

	now := time.Now()
	e = db.SetTimeValue("mytimekey", now)
	if e != nil {
		t.Fatal(e)
	}

	stamp, e = db.GetTimeValue("mytimekey")
	if e != nil {
		t.Fatal(e)
	}
	if got, want := stamp, now; !want.Equal(want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestGetSetBytes(t *testing.T) {
	db, cleanup := newTestDB()
	defer cleanup()

	// Get non-existent key returns nil
	val, err := db.GetBytes("binkey")
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Fatalf("expected nil, got %v", val)
	}

	// Set and get bytes
	data := []byte{0x01, 0x02, 0x03, 0xFF}
	if err := db.SetBytes("binkey", data); err != nil {
		t.Fatal(err)
	}

	val, err = db.GetBytes("binkey")
	if err != nil {
		t.Fatal(err)
	}
	if len(val) != len(data) {
		t.Fatalf("got len %d, want %d", len(val), len(data))
	}
	for i := range data {
		if val[i] != data[i] {
			t.Fatalf("byte %d: got %x, want %x", i, val[i], data[i])
		}
	}
}

func TestDeleteBucket(t *testing.T) {
	db, cleanup := newTestDB()
	defer cleanup()

	// Delete non-existent bucket should not error
	if err := db.DeleteBucket("nonexistent"); err != nil {
		t.Fatalf("expected no error deleting nonexistent bucket, got %s", err)
	}

	// Create data in the awless bucket, then delete it
	if err := db.SetStringValue("key1", "val1"); err != nil {
		t.Fatal(err)
	}
	if err := db.DeleteBucket(awlessBucket); err != nil {
		t.Fatal(err)
	}

	// After deleting bucket, values should be gone (returns empty)
	val, err := db.GetStringValue("key1")
	if err != nil {
		t.Fatal(err)
	}
	if val != "" {
		t.Fatalf("expected empty string after bucket deletion, got %q", val)
	}
}

func TestExecute(t *testing.T) {
	f, e := os.MkdirTemp(".", "testexec")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(f)
	os.Setenv("__AWLESS_HOME", f)

	err := Execute(func(db *DB) error {
		if err := db.SetStringValue("execkey", "execval"); err != nil {
			return err
		}
		val, err := db.GetStringValue("execkey")
		if err != nil {
			return err
		}
		if val != "execval" {
			t.Fatalf("got %q, want %q", val, "execval")
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecuteNoHome(t *testing.T) {
	os.Setenv("__AWLESS_HOME", "")
	err := Execute(func(db *DB) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error when __AWLESS_HOME is not set")
	}
}

func TestCloseNilBolt(t *testing.T) {
	// Close should not panic when bolt is nil
	db := &DB{bolt: nil}
	db.Close() // should not panic
}

func TestOpenInvalidPath(t *testing.T) {
	_, err := open("/nonexistent/path/to/db")
	if err == nil {
		t.Fatal("expected error opening db at invalid path")
	}
}

func TestSetAndGetTimeValueRoundTrip(t *testing.T) {
	db, cleanup := newTestDB()
	defer cleanup()

	// Test with a specific time to ensure proper round-trip
	fixedTime := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	if err := db.SetTimeValue("fixedtime", fixedTime); err != nil {
		t.Fatal(err)
	}
	got, err := db.GetTimeValue("fixedtime")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Equal(fixedTime) {
		t.Fatalf("got %s, want %s", got, fixedTime)
	}
}

func TestGetIntValueNonNumeric(t *testing.T) {
	db, cleanup := newTestDB()
	defer cleanup()

	// Store a non-numeric string, then try to get as int
	if err := db.SetStringValue("badint", "notanumber"); err != nil {
		t.Fatal(err)
	}
	_, err := db.GetIntValue("badint")
	if err == nil {
		t.Fatal("expected error converting non-numeric string to int")
	}
}

func TestSetBytesOverwrite(t *testing.T) {
	db, cleanup := newTestDB()
	defer cleanup()

	if err := db.SetBytes("key", []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := db.SetBytes("key", []byte("second")); err != nil {
		t.Fatal(err)
	}
	val, err := db.GetBytes("key")
	if err != nil {
		t.Fatal(err)
	}
	if string(val) != "second" {
		t.Fatalf("got %q, want %q", string(val), "second")
	}
}
