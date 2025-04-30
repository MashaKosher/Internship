package entity

type SearchParams struct {
	Name string `json:"name" validate:"omitempty,min=1,max=100,name_format"`
}
