package schema

// Tenant ...
type Tenant struct {
	ID     string `json:"id,omitempty"`
	Status int    `json:"status"`
	OTID   string `json:"tenant"`
}

// Subject ...
type Subject struct {
	ID      string `json:"id,omitempty"`
	Status  int    `json:"status"`
	Subject string `json:"subject"`
}

// SubjectUnit ...
type SubjectUnit struct {
	*Unit
	Level int    `json:"level"`
	Kind  string `json:"kind"`
	Name  string `json:"name"`
}

// Unit ...
type Unit struct {
	ID         string `json:"id,omitempty"`
	Status     int    `json:"status"`
	TargetID   string `json:"targetId"`
	TargetType string `json:"targetType"`
}

// Object ...
type Object struct {
	ID         string `json:"id,omitempty"`
	TargetID   string `json:"targetId"`
	TargetType string `json:"targetType"`
	Search     string `json:"search,omitempty"`
}

// Permission ...
type Permission struct {
	ID         string `json:"id,omitempty"`
	Permission string `json:"permission"`
}

// PermissionWithFacets ...
type PermissionWithFacets struct {
	*Permission
	Extention string `json:"extention"`
}

// Scope ...
type Scope struct {
	ID         string `json:"id,omitempty"`
	Status     int    `json:"status"`
	TargetID   string `json:"targetId"`
	TargetType string `json:"targetType"`
}
