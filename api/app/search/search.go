package search

type OrderBy string

const (
	OrderByAlphabetical	OrderBy =	"ALPHABETICAL"
	OrderBySortOrder	OrderBy =	"SORT_ORDER"
	OrderByCreatedAt	OrderBy =	"CREATED_AT"
	OrderByUpdatedAt	OrderBy =	"UPDATED_AT"
)

var validOrderByValues = map[OrderBy]struct{}{
	OrderByAlphabetical:	{},
	OrderBySortOrder:	{},
	OrderByCreatedAt:	{},
	OrderByUpdatedAt:	{},
}

func IsValidOrderBy(orderBy string) bool {
	_, valid := validOrderByValues[OrderBy(orderBy)]
	return valid
}