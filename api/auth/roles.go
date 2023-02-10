// Package auth contains the rule of authentication and authorization used by the application
package auth

// Roles available
const (
	ManagerRole    = uint32(1 << iota) // 1
	TechnicianRole                     // 2
)
