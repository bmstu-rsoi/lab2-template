package ratingdb

import (
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model

	Username string `gorm:"size:80;uniqueIndex;not null;<-:create"`
	Stars    uint32 `gorm:"check:stars BETWEEN 0 AND 100;default:0"`
}
