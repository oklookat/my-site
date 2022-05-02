package base

// get paginated articles by params.
type ArticleGetParams struct {
	// number of page.
	Page int

	// show published or drafts?
	Drafts bool

	// start from newest? true == DESC; false == ASC.
	Newest bool

	// created; updated; published.
	By string

	// search by title.
	Title *string
}

// article request body that user should send to create/update article.
type ArticleBody struct {
	CoverID     *string `json:"cover_id"`
	IsPublished *bool   `json:"is_published"`
	Title       *string `json:"title"`
	Content     *string `json:"content"`
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
