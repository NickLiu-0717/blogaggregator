package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	config "github.com/NickLiu-0717/blogaggregator/internal/config"
	"github.com/NickLiu-0717/blogaggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Printf("Couldn't read config: %s", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Printf("Couldn't open postgres from dburl: %s", err)
	}

	dbQueries := database.New(db)

	cfgstate := state{
		cfg: &cfg,
		db:  dbQueries,
	}
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerlogin)
	cmds.register("register", handlerregister)

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}

	cmd := command{
		name:      os.Args[1],
		arguments: os.Args,
	}
	err = cmds.run(&cfgstate, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
