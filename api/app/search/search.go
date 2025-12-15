package search

type OrderBy string

const (
	OrderBySortOrder	OrderBy =	"SORT_ORDER"
	OrderByAlphabetical	OrderBy =	"ALPHABETICAL"
	OrderByCreatedAt	OrderBy =	"CREATED_AT"
	OrderByUpdatedAt	OrderBy =	"UPDATED_AT"
)

var validOrderByValues = map[OrderBy]struct{}{
	OrderBySortOrder:	{},
	OrderByAlphabetical:	{},
	OrderByCreatedAt:	{},
	OrderByUpdatedAt:	{},
}

func IsValidOrderBy(orderBy string) bool {
	_, valid := validOrderByValues[OrderBy(orderBy)]
	return valid
}