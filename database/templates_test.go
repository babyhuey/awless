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
	"encoding/json"
	"testing"

	"github.com/wallix/awless/template"
)

func newTestTemplateExecution(id string, cmdLine string) *template.TemplateExecution {
	tpl, err := template.Parse(cmdLine)
	if err != nil {
		panic(err)
	}
	tpl.ID = id
	return &template.TemplateExecution{
		Template: tpl,
		Author:   "testauthor",
		Source:   "testsource",
		Locale:   "us-east-1",
		Profile:  "default",
	}
}

func TestAddAndGetTemplate(t *testing.T) {
	db, cleanup := newTestDb()
	defer cleanup()

	tplExec := newTestTemplateExecution("tpl-001", "create instance name=test")

	if err := db.AddTemplate(tplExec); err != nil {
		t.Fatal(err)
	}

	got, err := db.GetTemplate("tpl-001")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "tpl-001" {
		t.Fatalf("got ID %q, want %q", got.ID, "tpl-001")
	}
	if got.Author != "testauthor" {
		t.Fatalf("got Author %q, want %q", got.Author, "testauthor")
	}
	if got.Source != "testsource" {
		t.Fatalf("got Source %q, want %q", got.Source, "testsource")
	}
	if got.Locale != "us-east-1" {
		t.Fatalf("got Locale %q, want %q", got.Locale, "us-east-1")
	}
	if got.Profile != "default" {
		t.Fatalf("got Profile %q, want %q", got.Profile, "default")
	}
}

func TestAddTemplateEmptyID(t *testing.T) {
	db, cleanup := newTestDb()
	defer cleanup()

	tplExec := newTestTemplateExecution("", "create instance name=test")
	tplExec.ID = ""

	err := db.AddTemplate(tplExec)
	if err == nil {
		t.Fatal("expected error for empty ID, got nil")
	}
}

func TestGetTemplateNotFound(t *testing.T) {
	db, cleanup := newTestDb()
	defer cleanup()

	_, err := db.GetTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent template, got nil")
	}
}

func TestGetTemplateNoBucket(t *testing.T) {
	db, cleanup := newTestDb()
	defer cleanup()

	// No templates added yet, so bucket doesn't exist
	_, err := db.GetTemplate("any-id")
	if err == nil {
		t.Fatal("expected error when no templates bucket exists")
	}
}

func TestDeleteTemplates(t *testing.T) {
	db, cleanup := newTestDb()
	defer cleanup()

	// Delete when no bucket exists should succeed
	if err := db.DeleteTemplates(); err != nil {
		t.Fatal(err)
	}

	tplExec1 := newTestTemplateExecution("tpl-001", "create instance name=test1")
	tplExec2 := newTestTemplateExecution("tpl-002", "create instance name=test2")

	if err := db.AddTemplate(tplExec1); err != nil {
		t.Fatal(err)
	}
	if err := db.AddTemplate(tplExec2); err != nil {
		t.Fatal(err)
	}

	// Verify both exist
	if _, err := db.GetTemplate("tpl-001"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.GetTemplate("tpl-002"); err != nil {
		t.Fatal(err)
	}

	// Delete all
	if err := db.DeleteTemplates(); err != nil {
		t.Fatal(err)
	}

	// Verify all gone
	_, err := db.GetTemplate("tpl-001")
	if err == nil {
		t.Fatal("expected error after deleting all templates")
	}
}

func TestDeleteTemplate(t *testing.T) {
	db, cleanup := newTestDb()
	defer cleanup()

	// Delete when no bucket exists should error
	err := db.DeleteTemplate("tpl-001")
	if err == nil {
		t.Fatal("expected error when deleting from nonexistent bucket")
	}

	tplExec1 := newTestTemplateExecution("tpl-001", "create instance name=test1")
	tplExec2 := newTestTemplateExecution("tpl-002", "create instance name=test2")

	if err := db.AddTemplate(tplExec1); err != nil {
		t.Fatal(err)
	}
	if err := db.AddTemplate(tplExec2); err != nil {
		t.Fatal(err)
	}

	// Delete one
	if err := db.DeleteTemplate("tpl-001"); err != nil {
		t.Fatal(err)
	}

	// tpl-001 should be gone (Get returns empty content error)
	_, err = db.GetTemplate("tpl-001")
	if err == nil {
		t.Fatal("expected error for deleted template")
	}

	// tpl-002 should still exist
	got, err := db.GetTemplate("tpl-002")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "tpl-002" {
		t.Fatalf("got ID %q, want %q", got.ID, "tpl-002")
	}
}

func TestGetLoadedTemplate(t *testing.T) {
	db, cleanup := newTestDb()
	defer cleanup()

	// Error when no bucket
	_, err := db.GetLoadedTemplate("any-id")
	if err == nil {
		t.Fatal("expected error when no templates bucket exists")
	}

	tplExec := newTestTemplateExecution("tpl-001", "create instance name=test")

	if err := db.AddTemplate(tplExec); err != nil {
		t.Fatal(err)
	}

	loaded, err := db.GetLoadedTemplate("tpl-001")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Key != "tpl-001" {
		t.Fatalf("got Key %q, want %q", loaded.Key, "tpl-001")
	}
	if loaded.TplExec == nil {
		t.Fatal("expected TplExec to be non-nil")
	}
	if loaded.TplExec.ID != "tpl-001" {
		t.Fatalf("got TplExec.ID %q, want %q", loaded.TplExec.ID, "tpl-001")
	}
	if loaded.Raw == "" {
		t.Fatal("expected Raw to be non-empty")
	}
	// Verify Raw is valid JSON
	var js json.RawMessage
	if err := json.Unmarshal([]byte(loaded.Raw), &js); err != nil {
		t.Fatalf("Raw is not valid JSON: %s", err)
	}

	// Nonexistent ID with existing bucket
	_, err = db.GetLoadedTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent template ID")
	}
}

func TestListTemplates(t *testing.T) {
	db, cleanup := newTestDb()
	defer cleanup()

	// Empty list when no bucket
	results, err := db.ListTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 templates, got %d", len(results))
	}

	tplExec1 := newTestTemplateExecution("tpl-001", "create instance name=test1")
	tplExec2 := newTestTemplateExecution("tpl-002", "create instance name=test2")
	tplExec3 := newTestTemplateExecution("tpl-003", "create subnet name=test3")

	if err := db.AddTemplate(tplExec1); err != nil {
		t.Fatal(err)
	}
	if err := db.AddTemplate(tplExec2); err != nil {
		t.Fatal(err)
	}
	if err := db.AddTemplate(tplExec3); err != nil {
		t.Fatal(err)
	}

	results, err = db.ListTemplates()
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 templates, got %d", len(results))
	}

	// Verify each loaded template has key and raw content
	for _, lt := range results {
		if lt.Key == "" {
			t.Fatal("expected Key to be non-empty")
		}
		if lt.Raw == "" {
			t.Fatal("expected Raw to be non-empty")
		}
		if lt.TplExec == nil {
			t.Fatal("expected TplExec to be non-nil")
		}
	}
}
