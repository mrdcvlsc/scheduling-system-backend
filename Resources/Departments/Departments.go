package Departments

type Department struct {
	DepartmentID         uint16 `json:"DepartmentID,omitempty" bson:"DepartmentID,omitempty"`
	Code                 string `json:"Code" bson:"Code"`
	Name                 string `json:"Name,omitempty" bson:"Name,omitempty"`
	SaltedHashedPassword string `json:"SaltedHashedPassword,omitempty" bson:"SaltedHashedPassword,omitempty"`
}
