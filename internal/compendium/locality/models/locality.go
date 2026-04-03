package models

import "gesitr/internal/shared"

type Locality struct {
	shared.BaseModel `tstype:",extends"`
	Name             string `json:"name"`
	Owner            string `json:"owner"`
	Public           bool   `json:"public"`
}
