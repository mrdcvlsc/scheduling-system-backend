package Departments

type Department struct {
	DepartmentID uint16 `json:"DepartmentID" bson:"DepartmentID"`
	Code         string `json:"Code" bson:"Code"`
	Name         string `json:"Name" bson:"Name"`
}
