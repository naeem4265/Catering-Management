package data

// Restaurant Information
type Restaurant struct {
	Id       int    `json:"restaurantId"`
	Name     string `json:"restaurantName"`
	Location string `json:"restaurantLocation"`
}

// Menu Information
type Menu struct {
	Id     int    `json:"menuId"`
	RestId int    `json:"restaurantId"`
	Name   string `json:"menuName"`
	Price  int    `json:"menuPrice"`
	Vote   int    `json:"menuVote"`
}
