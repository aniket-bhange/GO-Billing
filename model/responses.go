package model

type ClientResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	User      Users  `json:"user"`
	ClientRef Client `json:"client"`
}
