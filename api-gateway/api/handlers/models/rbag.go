package models

type Policy struct{
	User string `json:"user"`
	Domain string `json:"domain"`
	Action string `json:"action"`
}

type Role struct{
	Id string `json:"id"`
	Type string `json:"type"`
}