package models

type RegisteredModel struct {
	Name           string `json:"name"`
	RecipeAddress  string `json:"recipe_address"`
	Nodes          int    `json:"nodes"`
}
