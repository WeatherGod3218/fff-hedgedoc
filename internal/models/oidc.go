package models

type KeycloakGroup struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	SubGroupCount int             `json:"subGroupCount"`
	SubGroups     []KeycloakGroup `json:"subGroups"`
}

type KeycloakUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
