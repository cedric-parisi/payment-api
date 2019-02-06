package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultLimit  = 100
	defaultOffset = 0
)

// Filter represents filtering option
type Filter struct {
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Sorting []Sort `json:"-"`
}

// Sort represents the sorting options
type Sort struct {
	Field      string
	Descending bool
}

// FilteredList represents a page of results
// Stores the filters applied and the resource name
type FilteredList struct {
	Filter
	Resource   string      `json:"-"`
	Results    interface{} `json:"results"`
	TotalCount int         `json:"total_count"`
}

// GetFilter extracts filtering options from the url
func GetFilter(params url.Values) *Filter {
	var limit int
	var offset int
	var err error

	if limit, err = strconv.Atoi(params.Get("limit")); err != nil {
		limit = defaultLimit
	}

	if offset, err = strconv.Atoi(params.Get("offset")); err != nil {
		offset = defaultOffset
	}

	return &Filter{
		Limit:   limit,
		Offset:  offset,
		Sorting: getSorting(params.Get("sort")),
	}
}

// Headers build headers Link for pagination
func (f FilteredList) Headers() http.Header {
	currentOffset := f.Filter.Offset
	remaining := f.TotalCount % f.Filter.Limit

	// compute first page link
	f.Filter.Offset = 0
	first := fmt.Sprintf(`</%s/%s>; rel="first"`, f.Resource, f.Filter.String())

	// compute previous page link
	prevOffset := 0
	if currentOffset > 0 {
		prevOffset = currentOffset - f.Filter.Limit
	}
	f.Filter.Offset = prevOffset
	prev := fmt.Sprintf(`</%s/%s>; rel="prev"`, f.Resource, f.Filter.String())

	//compute next page link
	nextOffset := f.TotalCount - remaining
	if currentOffset+f.Filter.Limit < f.TotalCount {
		nextOffset = currentOffset + f.Filter.Limit
	}
	f.Filter.Offset = nextOffset
	next := fmt.Sprintf(`</%s/%s>; rel="next"`, f.Resource, f.Filter.String())

	// compute last page link
	f.Filter.Offset = f.TotalCount - remaining
	if f.TotalCount%f.Limit == 0 {
		f.Filter.Offset = f.TotalCount - f.Limit
	}
	last := fmt.Sprintf(`</%s/%s>; rel="last"`, f.Resource, f.Filter.String())

	return http.Header{
		"Link": []string{first, prev, next, last},
	}
}

// String build a raw query according to the filters
func (f Filter) String() string {
	var sorts []string
	if len(f.Sorting) > 0 {
		for _, tmp := range f.Sorting {
			sorts = append(sorts, tmp.String())
		}
		return fmt.Sprintf("?limit=%d&offset=%d&sort=%s", f.Limit, f.Offset, strings.Join(sorts, ","))
	}
	return fmt.Sprintf("?limit=%d&offset=%d", f.Limit, f.Offset)
}

// String build the sort as part of a raw query
func (s Sort) String() string {
	sort := ""
	if s.Descending {
		sort += "-"
	}
	sort += s.Field
	return sort
}

func getSorting(source string) []Sort {
	strSorts := split(source)

	var sorting []Sort
	for _, sort := range strSorts {
		field := sort
		desc := false
		if sort[0] == '-' {
			field = sort[1:]
			desc = true
		}
		sorting = append(sorting, Sort{
			Field:      field,
			Descending: desc,
		})
	}
	return sorting
}

func split(source string) []string {
	// FieldsFunc split the source string each time the func(c rune) is true
	// Returns empty array if source is empty or every char statisfy func(c rune)
	// example : for source = ",,," it will return an empty array
	res := strings.FieldsFunc(source, func(c rune) bool {
		return c == ','
	})
	return res
}
