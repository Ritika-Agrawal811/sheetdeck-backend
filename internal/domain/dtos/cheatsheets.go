package dtos

type Cheatsheet struct {
	Slug        string `json:"slug" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Category    string `json:"category" binding:"required"`
	SubCategory string `json:"subcategory" binding:"required"`
}

type UpdateCheatsheetRequest struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	SubCategory string `json:"subcategory"`
}
