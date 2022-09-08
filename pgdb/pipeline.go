package pgdb

import (
	"time"

	"gorm.io/gorm"
)

type Pipeline struct {
	PipelineId 		string   	`gorm:"primaryKey"`
	DatasetId		string 		
	RegistryId		string
	ImageId			string
	SecretName 		string			
	Format 			string			
	Query 			string			
	HasHeader 		bool			
	Crontab 		time.Time		
	Output 			string			
	Cloudname 		string			
}

type Image struct {
	ImageId 		string		`gorm:"primaryKey"`
	Name 			string
	Version 		string
	Registry 		string
}

type Dataset struct {
	DatasetId 		string		`gorm:"primaryKey"`
	Local 			bool
	Path 			string
	Key 			string
}

type Registry struct {
	RegistryId		string		`gorm:"primaryKey"`
	Name 			string
	User 			string
	Password 		string
	Key 			string
}

func (p *Pipeline) Insert(db *gorm.DB) {
	err := db.Create(p)
	if err.Error != nil {
		panic(err.Error)
	}
}

func (r *Registry) Insert(db *gorm.DB) {
	err := db.Create(r)
	if err.Error != nil {
		panic(err.Error)
	}
}

func (i *Image) Insert(db *gorm.DB) {
	err := db.Create(i)
	if err.Error != nil {
		panic(err.Error)
	}
}


