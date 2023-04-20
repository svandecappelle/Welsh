package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStore struct {
	ingredient []*Ingredient
}

func (t testStore) Open() error {
	return nil
}

func (t testStore) Close() error {
	return nil
}

func (t testStore) GetIngredient() ([]*Ingredient, error) {
	return t.ingredient, nil
}

func (t testStore) GetRecipe() ([]*Recipe, error) {
	return nil, nil
}

func (t testStore) GetIngredientToRecipe(id int64) ([]*Ingredient, error) {
	return nil, nil
}

func (t testStore) GetFavorite(username string) ([]*Recipe, error) {
	return nil, nil
}

func (t testStore) QueryGet(data interface{}, query string) (interface{}, error) {
	return nil, nil
}

func (t testStore) QueryRow(data interface{}, query string) error {
	if strings.Contains(query, "ingredient") {
		return sql.ErrNoRows
	}
	if query == "SELECT password FROM user WHERE username='john'" {
		password := "$2a$14$ZyyJSmoV1gJI9xRRJDdWYOjudCbKx9dZYo/GK1xeyprl7hsH08J7u"
		v := reflect.ValueOf(password)
		reflect.ValueOf(data).Elem().Set(v)
		return nil
	}
	return nil
}

func (t testStore) QueryExec(inserId bool, iquery string) (int64, error) {
	return 0, nil
}

func TestIngredientCreateUnit(t *testing.T) {
	//Create server with test DB
	srv := newServer()
	srv.store = &testStore{}

	//Prepare JSON body
	ingredient := struct {
		Name string `json:"name"`
	}{
		Name: "chou-fleur, moutarde à l'ancienne, cheddar, bière brune, jambon blanc, oeuf, tomate, fromage de chèvre frais, oignons",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(ingredient)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/api/ingredient", &buf)
	w := httptest.NewRecorder()

	srv.handleIngredientCreate()(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIngredientCreateIntregation(t *testing.T) {
	//Create server with test DB
	srv := newServer()
	srv.store = &testStore{}

	//Prepare JSON body
	ingredient := struct {
		Name string `json:"name"`
	}{
		Name: "chou-fleur",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(ingredient)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/api/ingredient", &buf)
	r.Header.Set("username", "john")
	r.Header.Set("password", "test")
	w := httptest.NewRecorder()

	srv.serveHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
