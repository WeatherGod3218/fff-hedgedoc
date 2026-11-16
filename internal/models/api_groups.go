package models

type GetGroupsInput struct {
	Cursor *string `json:"cursor"`
	Prefix string  `json:"prefix"`
}

type GetGroupsPiece struct {
	GroupName string `json:"group_name"`
	UUID      string `json:"uuid"`
	Selected  bool   `json:"selected"`
}

type GetGroupsOutput struct {
	Groups []*GetGroupsPiece `json:"groups"`
	Cursor string            `json:"cursor"`
}
