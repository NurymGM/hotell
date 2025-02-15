package models

import "gorm.io/gorm"

type Room struct {
    gorm.Model
    Type        string  `json:"type"`
    Price       float64 `json:"price"`
    Info        string  `json:"info"`
    IsAvailable bool    `json:"is_available"`
    Image       string  `json:"image"`
}