package forms

type PostForm struct {
	Title     string `form:"title" json:"title" binding:"required"`
	Content   string `form:"content" json:"content"`
	Published bool   `form:"published" json:"published"`
	AuthorId  int    `form:"authorId" json:"authorId" binding:"required"`
}
