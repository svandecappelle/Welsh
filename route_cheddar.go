package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

type jsonIngredient struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type jsonRecipe struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Instruction string `json:"instruction"`
}
type jsonRecipeIngredient struct {
	jsonRecipe
	Ingredient []jsonIngredient
}
type jsonFavorite struct {
	Username string   `json:"username"`
	Recipe   []string `json:"recipe"`
	Flag     bool     `json:"flag"`
}

func (s *server) handleIngredientList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ingredient, err := s.store.GetIngredient()
		if err != nil {
			log.Printf("Cannot load ingredient. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = make([]jsonIngredient, len(ingredient))
		for i, ing := range ingredient {
			resp[i] = mapIngredientToJson(ing)
		}

		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleRecipeList() http.HandlerFunc {
	type request struct {
		Ingredient bool `json:"ingredient"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != io.EOF && err != nil {
			log.Printf("Cannot parse ingredient body. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		recipe, err := s.store.GetRecipe()
		if err != nil {
			log.Printf("Cannot load recipe. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}
		// List Recipes with ingredients
		if req.Ingredient {
			var resp = make([]jsonRecipeIngredient, len(recipe))
			for i, rec := range recipe {
				ingredient, err := s.store.GetIngredientToRecipe(rec.ID)
				if err != nil {
					log.Printf("Cannot load ingredient to recipe. err=%v\n", err)
					s.respond(w, r, nil, http.StatusInternalServerError)
					return
				}

				var resp_ing = make([]jsonIngredient, len(ingredient))
				for i, ing := range ingredient {
					resp_ing[i] = mapIngredientToJson(ing)
				}

				resp[i] = mapRecipeIngredientToJson(rec, resp_ing)
			}
			s.respond(w, r, resp, http.StatusOK)
			// List only recipes
		} else {
			var resp = make([]jsonRecipe, len(recipe))
			for i, rec := range recipe {
				resp[i] = mapRecipeToJson(rec)
			}
			s.respond(w, r, resp, http.StatusOK)
		}
	}
}

func (s *server) handleRecipeDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Cannot parse id to int. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		res, err := s.store.QueryGet(&Recipe{}, "SELECT * FROM recipe WHERE id="+strconv.FormatInt(id, 10))
		if err != nil {
			log.Printf("Cannot load recipe, err=%v", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		recipe := res.(*Recipe)
		ingredient, err := s.store.GetIngredientToRecipe(recipe.ID)
		if err != nil {
			log.Printf("Cannot load ingredient to recipe. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp_ing = make([]jsonIngredient, len(ingredient))
		for i, ing := range ingredient {
			resp_ing[i] = mapIngredientToJson(ing)
		}

		resp := mapRecipeIngredientToJson(recipe, resp_ing)

		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleFavoriteList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("username")
		favorite, err := s.store.GetFavorite(user)
		if err != nil {
			log.Printf("Cannot load favorite. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = make([]jsonRecipe, len(favorite))
		for i, recipe := range favorite {
			resp[i] = mapRecipeToJson(recipe)
		}

		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleIngredientCreate() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse ingredient body. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		sReq := regexp.MustCompile(`(\s*,\s*)+`).Split(req.Name, -1)
		var resp = make([]jsonIngredient, len(sReq))
		for i, ingredient := range sReq {
			// Create ingredient
			ingredientStore := &Ingredient{
				ID:   0,
				Name: ingredient,
			}

			err := s.store.QueryRow(&ingredient, "SELECT name FROM ingredient WHERE name='"+ingredient+"'")
			if err != sql.ErrNoRows {
				msg := fmt.Sprintf("Ingredients already exist: %v", ingredient)
				respErr := &jsonErr{
					Message: msg,
				}
				log.Print(msg)
				s.respond(w, r, respErr, http.StatusBadRequest)
				return
			}
			if err != nil {
				if err == sql.ErrNoRows {
					//Store the ingredient in the DB
					id, err := s.store.QueryExec(true, "INSERT INTO ingredient (name) VALUES ('"+ingredientStore.Name+"')")
					if err != nil {
						log.Printf("Cannot create ingredient in DB. err=%v", err)
						s.respond(w, r, nil, http.StatusInternalServerError)
						return
					}
					ingredientStore.ID = id
				} else {
					log.Printf("Cannot load ingredient in DB. err=%v", err)
					s.respond(w, r, nil, http.StatusInternalServerError)
					return
				}
			}
			resp[i] = mapIngredientToJson(ingredientStore)

		}
		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleRecipeCreate() http.HandlerFunc {
	type requestIngredient struct {
		Name string `json:"name"`
	}
	type request struct {
		Name        string              `json:"name"`
		Description string              `json:"description"`
		Instruction string              `json:"instruction"`
		Ingredient  []requestIngredient `json:"ingredient"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse recipe body. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		// Create a recipe
		recipeStore := &Recipe{
			ID:          0,
			Name:        req.Name,
			Description: req.Description,
			Instruction: req.Instruction,
		}

		// Store the recipe in the DB
		log.Print("INSERT INTO recipe (name, description, instruction) VALUES ('" + recipeStore.Name + "', '" + recipeStore.Description + "', '" + recipeStore.Instruction + "')")
		id, err := s.store.QueryExec(true, "INSERT INTO recipe (name, description, instruction) VALUES ('"+recipeStore.Name+"', '"+recipeStore.Description+"', '"+recipeStore.Instruction+"')")
		if err != nil {
			log.Printf("Cannot create recipe in DB. err=%v", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		recipeStore.ID = id
		var respIng = make([]jsonIngredient, len(req.Ingredient))
		for i, ing := range req.Ingredient {
			res, err := s.store.QueryGet(&Ingredient{}, "SELECT * FROM ingredient WHERE name='"+ing.Name+"'")
			if err != nil {
				log.Printf("Cannot load ingredient, err=%v", err)
				s.respond(w, r, nil, http.StatusInternalServerError)
				return
			}
			ingredientStore := res.(*Ingredient)

			// Store recipe_id and ingredient_if
			_, err = s.store.QueryExec(false, "INSERT INTO recipeingredient (recipe_id, ingredient_id) VALUES ("+strconv.FormatInt(recipeStore.ID, 10)+", "+strconv.FormatInt(ingredientStore.ID, 10)+")")
			if err != nil {
				log.Printf("Cannot create the link between recipe and ingredient in DB. err=%v", err)
				s.respond(w, r, nil, http.StatusInternalServerError)
				return
			}

			respIng[i] = mapIngredientToJson(ingredientStore)
		}

		var resp = mapRecipeIngredientToJson(recipeStore, respIng)
		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleFavoriteFlag() http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
		Flag bool   `json:"flag"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse recipe body. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		user := r.Header.Get("username")
		sReq := regexp.MustCompile(`(\s*,\s*)+`).Split(req.Name, -1)

		res, err := s.store.QueryGet(&User{}, "SELECT * FROM user WHERE username='"+user+"'")
		if err != nil {
			log.Printf("Cannot load user, err=%v", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		userStore := res.(*User)
		for _, recipe := range sReq {
			res, err := s.store.QueryGet(&Recipe{}, "SELECT * FROM recipe WHERE name='"+recipe+"'")
			if err != nil {
				log.Printf("Cannot load recipe, err=%v", err)
				s.respond(w, r, nil, http.StatusInternalServerError)
				return
			}

			recipeStore := res.(*Recipe)
			if req.Flag {
				_, err = s.store.QueryExec(false, "INSERT INTO favorite (user_id, recipe_id) VALUES ("+strconv.FormatInt(userStore.ID, 10)+", "+strconv.FormatInt(recipeStore.ID, 10)+")")
				if err != nil {
					log.Printf("Cannot create favorite in DB. err=%v", err)
					s.respond(w, r, nil, http.StatusInternalServerError)
					return
				}
				resp := mapFavoriteToJson(userStore, sReq, true)
				s.respond(w, r, resp, http.StatusOK)
			} else {
				_, err = s.store.QueryExec(false, "DELETE FROM favorite WHERE user_id='"+strconv.FormatInt(userStore.ID, 10)+"' AND recipe_id='"+strconv.FormatInt(recipeStore.ID, 10)+"'")
				if err != nil {
					log.Printf("Cannot delete favorite in DB. err=%v", err)
					s.respond(w, r, nil, http.StatusInternalServerError)
					return
				}
				resp := mapFavoriteToJson(userStore, sReq, false)
				s.respond(w, r, resp, http.StatusOK)
			}
		}
	}
}

func mapIngredientToJson(i *Ingredient) jsonIngredient {
	return jsonIngredient{
		ID:   i.ID,
		Name: i.Name,
	}
}

func mapRecipeToJson(r *Recipe) jsonRecipe {
	return jsonRecipe{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Instruction: r.Instruction,
	}
}

func mapRecipeIngredientToJson(r *Recipe, i []jsonIngredient) jsonRecipeIngredient {
	return jsonRecipeIngredient{
		jsonRecipe: jsonRecipe{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Instruction: r.Instruction,
		},
		Ingredient: i,
	}
}

func mapFavoriteToJson(u *User, rec []string, f bool) jsonFavorite {
	return jsonFavorite{
		Username: u.Username,
		Recipe:   rec,
		Flag:     f,
	}
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
