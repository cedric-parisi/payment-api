// +build !integration

package utils

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestGetFilter(t *testing.T) {
	type args struct {
		params url.Values
	}
	tests := []struct {
		name string
		args args
		want *Filter
	}{
		{
			name: "ok",
			args: args{
				params: url.Values{
					"sort":   {"date,-views"},
					"limit":  {"50"},
					"offset": {"100"},
				},
			},
			want: &Filter{
				Limit:  50,
				Offset: 100,
				Sorting: []Sort{
					{
						Field:      "date",
						Descending: false,
					},
					{
						Field:      "views",
						Descending: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFilter(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_String(t *testing.T) {
	type fields struct {
		Limit   int
		Offset  int
		Sorting []Sort
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "no sorting option ok",
			fields: fields{
				Limit:  100,
				Offset: 50,
			},
			want: "?limit=100&offset=50",
		},
		{
			name: "one sorting option ok",
			fields: fields{
				Limit:  100,
				Offset: 50,
				Sorting: []Sort{
					{
						Field:      "date",
						Descending: false,
					},
				},
			},
			want: "?limit=100&offset=50&sort=date",
		},
		{
			name: "multiple sorting options ok",
			fields: fields{
				Limit:  100,
				Offset: 50,
				Sorting: []Sort{
					{
						Field:      "date",
						Descending: false,
					},
					{
						Field:      "amount",
						Descending: true,
					},
				},
			},
			want: "?limit=100&offset=50&sort=date,-amount",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filter{
				Limit:   tt.fields.Limit,
				Offset:  tt.fields.Offset,
				Sorting: tt.fields.Sorting,
			}
			if got := f.String(); got != tt.want {
				t.Errorf("Filter.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilteredList_Headers(t *testing.T) {
	type fields struct {
		filteredList FilteredList
	}
	tests := []struct {
		name   string
		fields fields
		want   http.Header
	}{
		{
			name: "multiple page ok",
			fields: fields{
				filteredList: FilteredList{
					Resource:   "payments",
					TotalCount: 348,
					Filter: Filter{
						Limit:  50,
						Offset: 150,
						Sorting: []Sort{
							{
								Field:      "amount",
								Descending: true,
							},
						},
					},
				},
			},
			want: http.Header{
				"Link": []string{
					`</payments/?limit=50&offset=0&sort=-amount>; rel="first"`,
					`</payments/?limit=50&offset=100&sort=-amount>; rel="prev"`,
					`</payments/?limit=50&offset=200&sort=-amount>; rel="next"`,
					`</payments/?limit=50&offset=300&sort=-amount>; rel="last"`,
				},
			},
		},
		{
			name: "one page ok",
			fields: fields{
				filteredList: FilteredList{
					Resource:   "payments",
					TotalCount: 10,
					Filter: Filter{
						Limit:  50,
						Offset: 0,
						Sorting: []Sort{
							{
								Field:      "status",
								Descending: false,
							},
						},
					},
				},
			},
			want: http.Header{
				"Link": []string{
					`</payments/?limit=50&offset=0&sort=status>; rel="first"`,
					`</payments/?limit=50&offset=0&sort=status>; rel="prev"`,
					`</payments/?limit=50&offset=0&sort=status>; rel="next"`,
					`</payments/?limit=50&offset=0&sort=status>; rel="last"`,
				},
			},
		},
		{
			name: "total count is a multiplier of limit",
			fields: fields{
				filteredList: FilteredList{
					Resource:   "payments",
					TotalCount: 140,
					Filter: Filter{
						Limit:  10,
						Offset: 40,
						Sorting: []Sort{
							{
								Field:      "status",
								Descending: false,
							},
						},
					},
				},
			},
			want: http.Header{
				"Link": []string{
					`</payments/?limit=10&offset=0&sort=status>; rel="first"`,
					`</payments/?limit=10&offset=30&sort=status>; rel="prev"`,
					`</payments/?limit=10&offset=50&sort=status>; rel="next"`,
					`</payments/?limit=10&offset=130&sort=status>; rel="last"`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.filteredList
			if got := p.Headers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilteredList.Headers() = %v, want %v", got, tt.want)
			}
		})
	}
}
