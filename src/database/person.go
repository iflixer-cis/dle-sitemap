package database

import (
	"log"
	"path/filepath"
)

type Person struct {
	ID        int
	AltName   string
	Name      string
	URL       string
	UpdatedAt string
}

func (c *Person) TableName() string {
	return "dle_persons"
}

func (s *Service) PersonsAll() (persons []*Person, err error) {
	if err = s.DB.Find(&persons).Error; err != nil {
		log.Println("Cannot load persons", err)
	}
	for i, p := range persons {
		persons[i].URL = filepath.Join("person", p.AltName)
	}
	return
}
