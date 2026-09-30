package database

import (
	"log"
	"strconv"
)

type Person struct {
	ID        int
	Name      string
	URL       string
	UpdatedAt string
}

func (c *Person) TableName() string {
	return "flix_people"
}

func (s *Service) PeopleAll() (people []*Person, err error) {
	if err = s.DB.Find(&people).Error; err != nil {
		log.Println("Cannot load people", err)
	}
	for i, person := range people {
		people[i].URL = "people/" + strconv.Itoa(person.ID)
	}
	return
}
