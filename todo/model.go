package model

type Task struct {
	ID   int16 // huruf kecil dianggap private dong
	Name string
	Due  string
	Done bool
}
