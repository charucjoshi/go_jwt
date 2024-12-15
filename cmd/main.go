package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/charucjoshi/go_jwt/cmd/api"
	"github.com/charucjoshi/go_jwt/db"
	"github.com/charucjoshi/go_jwt/configs"
)

func main() {
	cfg := mysql.Config {
		User:                   configs.Envs.DBUser,
		Passwd:                 configs.Envs.DBPassword,
		Net:                    "tcp",
		Addr:                   configs.Envs.DBAddress,
		DBName:                 configs.Envs.DBName,
		AllowNativePasswords:   true,
		ParseTime:              true,
	}

	db, err := db.NewMySQLStorage(cfg)

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully connected!")
}
