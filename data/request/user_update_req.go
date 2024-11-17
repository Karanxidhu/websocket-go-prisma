package request


type UserUpdateReq struct {
	UserId string `json:"userId"`
	Username string `json:"username"`
}