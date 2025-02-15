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

type Inventory struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ReceivedCoins struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoins struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []ReceivedCoins `json:"received"`
	Sent     []SentCoins     `json:"sent"`
}

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}
