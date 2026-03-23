package shared

import (
	"testing"
)

func TestResolvePermissions_Owner(t *testing.T) {
	perms, visible := ResolvePermissions("user1", "user1", false)
	if !visible {
		t.Fatal("owner should be visible")
	}
	if len(perms) != 3 {
		t.Fatalf("owner should have 3 permissions, got %d", len(perms))
	}
	expected := map[Permission]bool{PermissionRead: true, PermissionModify: true, PermissionDelete: true}
	for _, p := range perms {
		if !expected[p] {
			t.Errorf("unexpected permission: %s", p)
		}
	}
}

func TestResolvePermissions_OwnerPublic(t *testing.T) {
	perms, visible := ResolvePermissions("user1", "user1", true)
	if !visible {
		t.Fatal("owner should be visible")
	}
	if len(perms) != 3 {
		t.Fatalf("owner of public entity should still have 3 permissions, got %d", len(perms))
	}
}

func TestResolvePermissions_NonOwnerPublic(t *testing.T) {
	perms, visible := ResolvePermissions("user2", "user1", true)
	if !visible {
		t.Fatal("public entity should be visible to non-owner")
	}
	if len(perms) != 1 || perms[0] != PermissionRead {
		t.Fatalf("non-owner on public entity should get [READ], got %v", perms)
	}
}

func TestResolvePermissions_NonOwnerPrivate(t *testing.T) {
	perms, visible := ResolvePermissions("user2", "user1", false)
	if visible {
		t.Fatal("private entity should not be visible to non-owner")
	}
	if perms != nil {
		t.Fatalf("non-owner on private entity should get nil, got %v", perms)
	}
}
