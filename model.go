package main

type Ingredient struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type Recipe struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Instruction string `db:"instruction"`
}

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
