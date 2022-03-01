package base

// get paginated articles by params.
type ArticleGetParams struct {
	// number of page.
	Page int
	// published; drafts.
	Show string
	// created; updated; published.
	By string
	// newest (DESC); oldest (ASC).
	Start string
	// true (with content); false (gives you empty content).
	Preview bool
	// category name.
	CategoryName *string
	// show articles only without category
	WithoutCategory bool
}

// article request body that user should send to create/update article.
type ArticleBody struct {
	CategoryID  *string `json:"category_id"`
	CoverID     *string `json:"cover_id"`
	IsPublished *bool   `json:"is_published"`
	Title       *string `json:"title"`
	Content     *string `json:"content"`
}

// category request body that user should send to create/update category.
type CategoryBody struct {
	// category name.
	Name string `json:"name"`
}

// get paginated files by params.
type FileGetParams struct {
	// number of page.
	Page int
	// newest (DESC); oldest (ASC).
	Start string
	// created (by creation date).
	By string
	// file extension without dot.
	Extension *string
	// OR extension type. image || audio || video.
	ExtensionType *string
}
