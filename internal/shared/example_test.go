package shared

import "fmt"

func ExampleResolvePermissions_owner() {
	perms, visible := ResolvePermissions("alice", "alice", false)
	fmt.Println(perms, visible)
	// Output: [READ MODIFY DELETE] true
}

func ExampleResolvePermissions_nonOwnerPublic() {
	perms, visible := ResolvePermissions("bob", "alice", true)
	fmt.Println(perms, visible)
	// Output: [READ] true
}

func ExampleResolvePermissions_nonOwnerPrivate() {
	perms, visible := ResolvePermissions("bob", "alice", false)
	fmt.Println(perms, visible)
	// Output: [] false
}
