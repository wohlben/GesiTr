package handlers

import (
	"gesitr/internal/compendium/locality/models"
)

type ListLocalityAvailabilitiesInput struct {
	LocalityID  string `query:"localityId" doc:"Filter by locality ID"`
	EquipmentID string `query:"equipmentId" doc:"Filter by equipment ID"`
	Available   string `query:"available" doc:"Filter by availability ('true' or 'false')"`
}

type ListLocalityAvailabilitiesOutput struct {
	Body []models.LocalityAvailability
}

type CreateLocalityAvailabilityInput struct {
	Body struct {
		LocalityID  uint  `json:"localityId" required:"true"`
		EquipmentID uint  `json:"equipmentId" required:"true"`
		Available   *bool `json:"available,omitempty"`
	}
}

type CreateLocalityAvailabilityOutput struct {
	Body models.LocalityAvailability
}

type UpdateLocalityAvailabilityInput struct {
	ID   uint `path:"id"`
	Body struct {
		Available bool `json:"available" required:"true"`
	}
}

type UpdateLocalityAvailabilityOutput struct {
	Body models.LocalityAvailability
}

type DeleteLocalityAvailabilityInput struct {
	ID uint `path:"id"`
}

type DeleteLocalityAvailabilityOutput struct{}
