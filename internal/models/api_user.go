package models

type UserInfo struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}
