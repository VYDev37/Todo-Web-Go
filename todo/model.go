package model

type Task struct {
	ID   int16  `gorm:"primaryKey;autoIncrement" json:"ID"`
	Name string `gorm:"not null" json:"Name"`
	Due  string `gorm:"not null" json:"Due"`
	Done bool   `gorm:"default:false" json:"Done"`
}
