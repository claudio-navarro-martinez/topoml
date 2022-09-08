package pgdb

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
/* 	opts := &pg.Options{
		User: "postgres",
		Password: "mysecretpassword",
		Addr: "172.17.0.2:5432",
	}

	var db *pg.DB = pg.Connect(opts)
	if db == nil {
		log.Printf("Failed to connect to database")
	}
	return db */

	dbURL := "postgres://postgres:mysecretpassword@172.17.0.2:5432/postgres"

    db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

    if err != nil {
        log.Fatalln(err)
		fmt.Println("no se pudo abrir conexion")
    }
    return db
}