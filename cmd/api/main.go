package main

import (
	"flag"
	"log"

	"github.com/Fankemp/GameMatch/internal/config"
	"github.com/Fankemp/GameMatch/internal/db_conn"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "run database migrations")
	migrateFlagDown := flag.Bool("migrate-down", false, "rollback last migrations")
	flag.Parse()

	cfg := config.NewPostgreConfig()
	db, err := db_conn.NewDB(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if *migrateFlag {
		if err = db.Migrate(); err != nil {
			log.Fatalln(err)
		}
		log.Println("migrations applied successfully")
		return
	}

	if *migrateFlagDown {
		if err = db.MigrateDown(); err != nil {
			log.Println("migration cant be rollback")
			return
		}
	}
	
}
