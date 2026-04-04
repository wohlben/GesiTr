package models

import "gesitr/internal/shared"

type Locality struct {
	shared.BaseModel `tstype:",extends"`
	Name             string `json:"name"`
	OwnershipGroupID uint   `json:"ownershipGroupId"`
	Public           bool   `json:"public"`
}
