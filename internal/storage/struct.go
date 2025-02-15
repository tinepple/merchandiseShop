package storage

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type Item struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}

type Transaction struct {
	UserIDFrom int `db:"user_id_from"`
	UserIDTo   int `db:"user_id_to"`
	Amount     int `db:"amount"`
}

type Inventory struct {
	Name     string `db:"name"`
	Quantity int    `db:"quantity"`
}

type CoinsHistory struct {
	UserNameFrom string `db:"username_from"`
	UserIDFrom   int    `db:"user_id_from"`
	UserNameTo   string `db:"username_to"`
	UserIDTo     int    `db:"user_id_to"`
	Amount       int    `db:"amount"`
}
