package graphql

type allResourcesFilter struct {
	Q  string `json:"q"`
	ID string `json:"id"`
}

type allResourcesQuery struct {
	Page      int                `json:"page"`
	PerPage   int                `json:"perPage"`
	SortField string             `json:"sortField"`
	SortOrder string             `json:"sortOrder"`
	Filter    allResourcesFilter `json:"filter"`
}

type allResourcesMeta struct {
	Count int `json:"count"`
}
