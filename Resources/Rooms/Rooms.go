package rooms

type Room struct {
	RoomID       uint16
	DepartmentID uint16 `json:"DepartmentID"`
	Capacity     uint16 `json:"Capacity"`
	RoomType     uint32 `json:"RoomType"`
	Name         string `json:"Name"`
}
