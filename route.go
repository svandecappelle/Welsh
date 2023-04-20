package main

func (s *server) routes() {
	s.router.HandleFunc("/", s.HandleIndex()).Methods("GET")
	s.router.HandleFunc("/api/ingredient", s.loggedOnly(s.handleIngredientList())).Methods("GET")
	s.router.HandleFunc("/api/recipe", s.loggedOnly(s.handleRecipeList())).Methods("GET")
	s.router.HandleFunc("/api/recipe/{id:[0-9]+}", s.loggedOnly(s.handleRecipeDetail())).Methods("GET")
	s.router.HandleFunc("/api/favorite", s.loggedOnly(s.handleFavoriteList())).Methods("GET")
	s.router.HandleFunc("/api/ingredient", s.loggedOnly(s.handleIngredientCreate())).Methods("POST")
	s.router.HandleFunc("/api/recipe", s.loggedOnly(s.handleRecipeCreate())).Methods("POST")
	s.router.HandleFunc("/api/favorite", s.loggedOnly(s.handleFavoriteFlag())).Methods("POST")
	s.router.HandleFunc("/api/admin/user", s.handleCreateUser()).Methods("POST")
}
