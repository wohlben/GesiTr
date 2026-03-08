package models

type UserExercise struct {
	BaseModel            `tstype:",extends"`
	Owner                string `json:"owner"`
	CompendiumExerciseID string `json:"compendiumExerciseId"`
	CompendiumVersion    int    `json:"compendiumVersion"`
}
