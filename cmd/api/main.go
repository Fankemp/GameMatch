package main

import (
	"github.com/Fankemp/GameMatch/internal/config"
	"github.com/Fankemp/GameMatch/internal/db_conn"
)

func main() {
	cfg := config.NewPostgreConfig()
	db, err := db_conn.NewDB(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

}
