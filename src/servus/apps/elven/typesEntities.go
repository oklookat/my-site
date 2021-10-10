package elven

import "servus/core"

type entityBase struct {
	*core.BaseController
}

// entityUser - manage users.
type entityUser struct {
}

// entityToken - manage tokens.
type entityToken struct {
}

// entityArticle - manage articles.
type entityArticle struct {
	*entityBase
}

// entityFile - manage files.
type entityFile struct {
	*entityBase
}

// entityAuth - manage authorization.
type entityAuth struct {
	*entityBase
}

// objectCmd - commandline methods.
type objectCmd struct {
}

// objectUtils - useful utilities.
type objectUtils struct {
}
