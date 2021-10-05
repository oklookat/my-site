package elven

import (
	"github.com/gosimple/slug"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"servus/core"
	"time"
)

func dbArticleBeforeChangeHook(article *ModelArticle) {
	// create temp slug (ULID)
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	article.Slug = ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

func dbArticleAfterChangeHook(article *ModelArticle) (err error) {
	// create normal slug
	article.Slug = slug.Make(article.Title) + "-" + article.ID
	var sql = "UPDATE articles SET slug=:slug WHERE id=:id RETURNING *"
	row, err := core.Database.NamedQuery(sql, &article)
	if err != nil {
		return err
	}
	err = row.StructScan(&article)
	if err != nil {
		return err
	}
	return nil
}
