package model

type ApiResp struct {
	Code      int          `json:"code"`
	Token     string       `json:"token"`
	Expire    string       `json:"expire"`
	User      User         `json:"User"`
	UserRoles []*UserRoles `json:"UserRoles"`
}
type ApiRespMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
