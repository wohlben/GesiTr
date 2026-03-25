package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"gesitr/internal/database"
	"gesitr/internal/profile/models"
)

func TestGetMyProfile(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	// Create profile for test user
	database.DB.Create(&models.UserProfileEntity{ID: "alice", Name: "Alice"})

	r := newRouter()
	w := doJSON(r, "GET", "/api/user/profile", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var profile models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &profile)

	if profile.ID != "alice" {
		t.Errorf("expected id alice, got %s", profile.ID)
	}
	if profile.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", profile.Name)
	}
}

func TestGetMyProfile_NotFound(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	r := newRouter()
	w := doJSON(r, "GET", "/api/user/profile", nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateMyProfile(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	database.DB.Create(&models.UserProfileEntity{ID: "alice", Name: "Alice"})

	r := newRouter()
	w := doJSON(r, "PATCH", "/api/user/profile", models.UpdateProfileRequest{Name: "Alice Smith"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var profile models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &profile)

	if profile.Name != "Alice Smith" {
		t.Errorf("expected name Alice Smith, got %s", profile.Name)
	}
}

func TestUpdateMyProfile_BadRequest(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	database.DB.Create(&models.UserProfileEntity{ID: "alice", Name: "Alice"})

	r := newRouter()
	w := doJSON(r, "PATCH", "/api/user/profile", map[string]string{})

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestGetProfile(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	database.DB.Create(&models.UserProfileEntity{ID: "bob", Name: "Bob"})

	r := newRouter()
	w := doJSON(r, "GET", "/api/profiles/bob", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var profile models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &profile)

	if profile.ID != "bob" {
		t.Errorf("expected id bob, got %s", profile.ID)
	}
	if profile.Name != "Bob" {
		t.Errorf("expected name Bob, got %s", profile.Name)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	setupTestDB(t)
	defer closeDB(t)

	r := newRouter()
	w := doJSON(r, "GET", "/api/profiles/nonexistent", nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
