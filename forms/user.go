package forms

type UserForm struct {
	Email string `form:"email" json:"email" binding:"required"`
	Name  string `form:"name" json:"name" binding:"required"`
	Pass  string `form:"pass" json:"pass" binding:"required,min=3,max=50"`
	Role  string `form:"role" json:"role"`
}

type LoginForm struct {
	Email string `form:"email" json:"email" binding:"required,email"`
	Pass  string `form:"pass" json:"pass" binding:"required,min=3,max=50"`
}
