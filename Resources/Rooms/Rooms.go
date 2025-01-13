package rooms

type Room struct {
	RoomID       uint16
	DepartmentID uint16 `json:"DepartmentID"`
	Capacity     uint16 `json:"Capacity"` // Maximum numbers of section allocation in a room
	RoomType     uint16 `json:"RoomType"` // Determines the room types.
	Name         string `json:"Name"`
}
