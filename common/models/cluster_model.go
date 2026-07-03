package models

import (
	"time"

	"github.com/Riven-Spell/sparky/common/enum"
)

type ClusterModel struct {
	Nickname       string           `json:"nickname"`
	SparkrunRecipe string           `json:"sparkrun_recipe"`
	Status         enum.ModelStatus `json:"status"`
	Endpoint       *string          `json:"endpoint"`
	Since          time.Time        `json:"since"`
	NoList         bool             `json:"no_list"`
	NoEvict        bool             `json:"no_evict"`
}
