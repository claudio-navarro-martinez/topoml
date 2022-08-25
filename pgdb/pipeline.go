package pgdb

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Pipeline struct {
	Datasetname string
	Dockerimage string
	Dockerversion string
	Secretname string
	Format string
	Query string
	HasHeader bool
	Crontab time.Time
	Pipelinename string
	Output string
	Cloudname string
}

func (p *Pipeline) Insert(db *pg.DB) {
	_, err := db.Model(p).Insert()
	if err != nil {
		panic(err)
	}
}

// createSchema creates database schema for User and Story models.
func (p *Pipeline) CreateSchema(db *pg.DB) error {
    models := []interface{}{
        (*Pipeline)(nil),
    }

    for _, model := range models {
        err := db.Model(model).CreateTable(&orm.CreateTableOptions{})          
        
        if err != nil {
            return err
        }
		
    }
    return nil
}