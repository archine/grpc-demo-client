package dto

type CreateUserArg struct {
	Name string `json:"name" binding:"required,max=32"`
	Age  int    `json:"age" binding:"required,gte=0,lte=150"`
	// Email 正则表达式验证
	Email string `json:"email" binding:"required,email,max=64"`
}
