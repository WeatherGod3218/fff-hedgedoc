package models

type GetUsersInput struct {
	Cursor *string `json:"cursor"`
	Prefix string  `json:"prefix"`
}

type GetUsersPiece struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Selected bool   `json:"selected"`
}

type GetUsersOutput struct {
	Users  []*GetUsersPiece `json:"users"`
	Cursor string           `json:"cursor"`
}

type UserInfo struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}
