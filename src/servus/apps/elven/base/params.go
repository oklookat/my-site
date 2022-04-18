package base

// get paginated articles by params.
type ArticleGetParams struct {
	// number of page.
	Page int

	// show published or drafts?
	Published bool

	// start from newest? true == DESC; false == ASC.
	Newest bool

	// true = empty content; false = full content.
	Preview bool

	// show only without category?
	WithoutCategory bool

	// created; updated; published.
	By string

	// search by category name.
	CategoryName *string

	// search by title.
	Title *string
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

	// extensions without dot.
	Extensions []string

	// filename.
	Filename *string
}
