package pgdb

import (
	"log"

	pg "github.com/go-pg/pg/v10"
)

func Connect() *pg.DB {
	opts := &pg.Options{
		User: "postgres",
		Password: "mysecretpassword",
		Addr: "172.17.0.2:5432",
}
var db *pg.DB = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect to database")
	}
	return db
}