package main

import (
	"context"
	"fmt"
	"os"

	config "github.com/NickLiu-0717/blogaggregator/internal/config"
	"github.com/NickLiu-0717/blogaggregator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handlerfunc := c.handlers[cmd.name]
	if err := handlerfunc(s, cmd); err != nil {
		return err
	}
	return nil
}

func handlerlogin(s *state, cmd command) error {
	if cmd.arguments == nil || len(cmd.arguments) == 2 {
		return fmt.Errorf("a username is required")
	}
	username := cmd.arguments[2]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User %s login successfully\n", username)
	return nil
}

func handlerregister(s *state, cmd command) error {
	if cmd.arguments == nil || len(cmd.arguments) == 2 {
		return fmt.Errorf("a username is required")
	}
	username := cmd.arguments[2]
	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		fmt.Printf("username %s already exists", username)
		os.Exit(1)
	}
	dbUser, err := s.db.CreateUser(context.Background(), username)
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", dbUser)
	return nil
}
