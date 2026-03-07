package database

import (
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dir := t.TempDir()
		if err := InitDB(dir + "/test.db"); err != nil {
			t.Fatal(err)
		}
		if DB == nil {
			t.Fatal("DB should not be nil")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		err := InitDB("/proc/self/fd/-1/impossible.db")
		if err == nil {
			t.Error("expected error for invalid path")
		}
	})
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	Init()

	if DB == nil {
		t.Fatal("DB should not be nil after Init()")
	}
	if _, err := os.Stat("gesitr.db"); os.IsNotExist(err) {
		t.Error("gesitr.db file was not created")
	}
}
