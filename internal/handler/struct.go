package handler

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinsRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
