package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	Open() error
	Close() error

	GetIngredient() ([]*Ingredient, error)
	GetRecipe() ([]*Recipe, error)
	GetIngredientToRecipe(id int64) ([]*Ingredient, error)
	GetFavorite(username string) ([]*Recipe, error)
	QueryGet(data interface{}, query string) (interface{}, error)
	QueryRow(data interface{}, query string) error
	QueryExec(inserId bool, iquery string) (int64, error)
}

type dbStore struct {
	db *sqlx.DB
}

var schema = `
CREATE TABLE IF NOT EXISTS ingredient
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	UNIQUE(id),
	UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS recipe
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	description TEXT,
	instruction TEXT
);

CREATE TABLE IF NOT EXISTS recipeingredient
(
	recipe_id INTEGER NOT NULL,
	ingredient_id INTEGER NOT NULL,
	FOREIGN KEY(recipe_id) REFERENCES recipe(id),
	FOREIGN KEY(ingredient_id) REFERENCES ingredient(id)
);

create table if not exists user
(
	id INTEGER primary key autoincrement,
	username TEXT,
	password TEXT
);

CREATE TABLE IF NOT EXISTS favorite
(
	user_id INTEGER NOT NULL,
	recipe_id INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES user(id),
	FOREIGN KEY(recipe_id) REFERENCES recipe(id),
	UNIQUE(user_id, recipe_id)
);
`

// FOREIGN KEY(recipe_id)REFERENCES recipe(id),
//ingredient_id NOT NULL INTEGER FOREIGN KEY(ingredient_id)REFERENCES ingredient(id)

func (store *dbStore) Open() error {
	db, err := sqlx.Connect("sqlite3", "cheddar.db")
	if err != nil {
		return err
	}
	log.Println("Connected to DB")
	db.MustExec(schema)
	store.db = db
	return nil
}

func (store *dbStore) Close() error {
	return store.db.Close()
}

func (store *dbStore) GetIngredient() ([]*Ingredient, error) {
	var ingredient []*Ingredient
	err := store.db.Select(&ingredient, "SELECT * FROM ingredient")
	if err != nil {
		return ingredient, err
	}
	return ingredient, nil
}

func (store *dbStore) GetRecipe() ([]*Recipe, error) {
	var recipe []*Recipe
	err := store.db.Select(&recipe, "SELECT * FROM recipe")
	if err != nil {
		return recipe, err
	}
	return recipe, nil
}

func (store *dbStore) GetIngredientToRecipe(id int64) ([]*Ingredient, error) {
	var ingredient []*Ingredient
	err := store.db.Select(&ingredient, "SELECT i.id, i.name FROM recipe r JOIN recipeingredient ri on r.id = ri.recipe_id JOIN ingredient i on i.id = ri.ingredient_id WHERE r.id=$1", id)
	if err != nil {
		return ingredient, err
	}
	return ingredient, nil
}

func (store *dbStore) GetFavorite(username string) ([]*Recipe, error) {
	var recipe []*Recipe
	err := store.db.Select(&recipe, "SELECT r.id, r.name, r.description, r.instruction FROM recipe r JOIN favorite f on r.id = f.recipe_id JOIN user u on u.id = f.user_id WHERE u.username='john'", username)
	if err != nil {
		return recipe, err
	}

	return recipe, nil
}

func (store *dbStore) QueryGet(data interface{}, query string) (interface{}, error) {
	err := store.db.Get(data, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (store *dbStore) QueryRow(data interface{}, query string) error {
	err := store.db.QueryRow(query).Scan(data)
	if err != nil {
		return err
	}

	return nil
}

func (store *dbStore) QueryExec(insertId bool, query string) (int64, error) {
	res, err := store.db.Exec(query)
	if err != nil {
		return 0, err
	}

	if insertId {
		id, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}

		return id, nil
	}

	return 0, nil
}
