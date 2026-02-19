package models

type Nums struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

type User struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type UserWithAge struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserIsAdult struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	IsAdult bool   `json:"is_adult"`
}

type Feedback struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}
