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
	// fmt.Printf("Read config: %+v\n", cfg)

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
	cmds.register("reset", handlerreset)
	cmds.register("users", handlerlistusers)
	cmds.register("agg", handleraggregate)
	cmds.register("addfeed", middlewareLoggedIn(handleraddfeed))
	cmds.register("feeds", handlerfeeds)
	cmds.register("follow", middlewareLoggedIn(handlerfollow))
	cmds.register("following", middlewareLoggedIn(handlerfollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerbrowse))

	if len(os.Args) < 2 {
		fmt.Println("Need command!")
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
