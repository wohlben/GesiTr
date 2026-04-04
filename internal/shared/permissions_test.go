package shared

import (
	"testing"
)

// mockAccess implements AccessChecker for testing.
type mockAccess struct {
	read, modify, delete bool
}

func (m mockAccess) CanRead() bool   { return m.read }
func (m mockAccess) CanModify() bool { return m.modify }
func (m mockAccess) CanDelete() bool { return m.delete }

func TestResolvePermissionsFromAccess_FullAccess(t *testing.T) {
	perms, visible := ResolvePermissionsFromAccess(mockAccess{read: true, modify: true, delete: true}, false)
	if !visible {
		t.Fatal("should be visible")
	}
	if len(perms) != 3 {
		t.Fatalf("full access should have 3 permissions, got %d", len(perms))
	}
}

func TestResolvePermissionsFromAccess_ModifyOnly(t *testing.T) {
	perms, visible := ResolvePermissionsFromAccess(mockAccess{read: true, modify: true, delete: false}, false)
	if !visible {
		t.Fatal("should be visible")
	}
	if len(perms) != 2 {
		t.Fatalf("modify access should have 2 permissions, got %d", len(perms))
	}
	for _, p := range perms {
		if p == PermissionDelete {
			t.Error("should not have DELETE permission")
		}
	}
}

func TestResolvePermissionsFromAccess_ReadOnly(t *testing.T) {
	perms, visible := ResolvePermissionsFromAccess(mockAccess{read: true, modify: false, delete: false}, false)
	if !visible {
		t.Fatal("should be visible")
	}
	if len(perms) != 1 || perms[0] != PermissionRead {
		t.Fatalf("read-only access should get [READ], got %v", perms)
	}
}

func TestResolvePermissionsFromAccess_NoAccessPublic(t *testing.T) {
	perms, visible := ResolvePermissionsFromAccess(mockAccess{}, true)
	if !visible {
		t.Fatal("public entity should be visible")
	}
	if len(perms) != 1 || perms[0] != PermissionRead {
		t.Fatalf("no access on public entity should get [READ], got %v", perms)
	}
}

func TestResolvePermissionsFromAccess_NoAccessPrivate(t *testing.T) {
	perms, visible := ResolvePermissionsFromAccess(mockAccess{}, false)
	if visible {
		t.Fatal("private entity with no access should not be visible")
	}
	if perms != nil {
		t.Fatalf("should get nil permissions, got %v", perms)
	}
}

func TestResolvePermissionsFromAccess_FullAccessPublic(t *testing.T) {
	perms, visible := ResolvePermissionsFromAccess(mockAccess{read: true, modify: true, delete: true}, true)
	if !visible {
		t.Fatal("should be visible")
	}
	if len(perms) != 3 {
		t.Fatalf("full access on public entity should still have 3 permissions, got %d", len(perms))
	}
}
