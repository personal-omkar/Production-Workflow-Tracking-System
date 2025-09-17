package model

type RolePermission struct {
	RoleId       int `json:"RoleId" gorm:"role_id"`
	PermissionId int `json:"PermissionId" gorm:"permission_id"`
}
