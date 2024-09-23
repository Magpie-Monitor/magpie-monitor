package routing

import (
	"net/url"
)

func LookupQueryParam(query url.Values, key string) (string, bool) {
	return query.Get(key), query.Has(key)

}
