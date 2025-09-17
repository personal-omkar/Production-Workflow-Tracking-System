package model

type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// UserDetails
type UserDetails struct {
	UserName  string
	FirstName string
	LastName  string
	// Role      int //TODO- decide role
}
