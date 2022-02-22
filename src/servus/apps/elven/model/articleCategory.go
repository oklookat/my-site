package model

import (
	"time"
)

// represents category in database.
type ArticleCategory struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// get all categories.
func (a *ArticleCategory) GetAll() (cats map[int]*ArticleCategory, err error) {
	var query = "SELECT * FROM article_categories ORDER BY name ASC"
	cats, err = articleCatAdapter.GetRows(query)
	return
}

// create category.
func (a *ArticleCategory) Create() (err error) {
	var query = "INSERT INTO article_categories (name) VALUES ($1) RETURNING *"
	err = articleCatAdapter.Get(a, query, a.Name)
	return
}

// change name.
func (a *ArticleCategory) ChangeNameByID() (err error) {
	var query = "UPDATE article_categories SET name=$1 WHERE id=$2 RETURNING *"
	err = articleCatAdapter.Get(a, query, a.Name, a.ID)
	return
}

// find by id.
func (a *ArticleCategory) FindByID() (found bool, err error) {
	found = false
	var query = "SELECT * FROM article_categories WHERE id=$1 LIMIT 1"
	founded, err := articleCatAdapter.Find(query, a.ID)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*a = *founded
	}
	return
}

// find by name.
func (a *ArticleCategory) FindByName() (found bool, err error) {
	found = false
	var query = "SELECT * FROM article_categories WHERE name=$1 LIMIT 1"
	founded, err := articleCatAdapter.Find(query, a.Name)
	if err != nil {
		return
	}
	if founded != nil {
		found = true
		*a = *founded
	}
	return
}

// delete by id.
func (a *ArticleCategory) DeleteByID() (err error) {
	var query = "DELETE FROM article_categories WHERE id=$1"
	_, err = articleCatAdapter.Exec(query, a.ID)
	return
}

// delete by name.
func (a *ArticleCategory) DeleteByName() (err error) {
	var query = "DELETE FROM article_categories WHERE name=$1"
	_, err = articleCatAdapter.Exec(query, a.Name)
	return
}
