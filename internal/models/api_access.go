package models

type ModifyAccessInput struct {
	NoteID         string `json:"note_id"`
	IdentifierUUID string `json:"identifier_uuid"`
	MemberType     string `json:"member_type"` //user or group, groups take priority in blacklists, users in whitelist
}
