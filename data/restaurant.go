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

// Menu Item
type Item struct {
	Id       int    `json:"itemId"`
	Name     string `json:"itemName"`
	ResName  string `json:"resName,omitempty"`
	ResId    int    `json:"resId,omitempty"`
	Location string `json:"resLocation"`
	Price    int    `json:"itemPrice"`
	Vote     int    `json:"itemVote"`
	Date     string `json:"winDate,omitempty"`
}

// Daily Winner
type Daily_Winner struct {
	Date               string `json:"winDate"`
	WinnerManuId       int    `json:"winnerManuId"`
	WinnerRestaurantId int    `json:"winnerResId"`
}
