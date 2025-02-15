package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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
	_, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("unresgist command")
	}

	handlerfunc := c.handlers[cmd.name]
	if err := handlerfunc(s, cmd); err != nil {
		return err
	}
	return nil
}

func handlerlogin(s *state, cmd command) error {
	if len(cmd.arguments) != 3 {
		return fmt.Errorf("incorrect arguments for login, need login USERNAME")
	}
	username := cmd.arguments[2]
	_, err := s.db.GetUserFromName(context.Background(), username)
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
	if len(cmd.arguments) != 3 {
		return fmt.Errorf("incorrect arguments for register, need register USERNAME")
	}
	username := cmd.arguments[2]
	_, err := s.db.GetUserFromName(context.Background(), username)
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

func handlerreset(s *state, cmd command) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("incorrect argument for reset, need reset only")
	}
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Reset successfully")
	return nil
}

func handlerlistusers(s *state, cmd command) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("incorrect argument for listing user, need list only")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handleraggregate(s *state, cmd command) error {
	if len(cmd.arguments) != 3 {
		return fmt.Errorf("incorrect argument for agg, need agg TIME")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[2])
	if err != nil {
		log.Printf("Error parsing duration: %s", err)
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for range ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func handleraddfeed(s *state, cmd command, dbUser database.User) error {
	if cmd.arguments == nil || len(cmd.arguments) != 4 {
		return fmt.Errorf("incorrect arguments for adding feed, need addfeed FEEDNAME URL")
	}

	dbFeed, err := s.db.CreateFeed(context.TODO(), database.CreateFeedParams{
		Name:   cmd.arguments[2],
		Url:    cmd.arguments[3],
		UserID: dbUser.ID,
	})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.TODO(), database.CreateFeedFollowParams{
		FeedID: dbFeed.ID,
		UserID: dbUser.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", dbFeed)
	return nil
}

func handlerfeeds(s *state, cmd command) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("incorrect argument for feeds, need feeds only")
	}
	feeds, err := s.db.GetFeeds(context.TODO())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		usrName, err := s.db.GetUserNameFromID(context.TODO(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed ID: %v, Feed name: %s, Feed URL: %s, User: %s\n", feed.ID, feed.Name, feed.Url, usrName)
	}
	return nil
}

func handlerfollow(s *state, cmd command, dbUser database.User) error {
	if len(cmd.arguments) != 3 {
		return fmt.Errorf("incorrect argument for follow, need follow URL")
	}
	dbFeed, err := s.db.GetFeedIDandNameFromURL(context.TODO(), cmd.arguments[2])
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.TODO(), database.CreateFeedFollowParams{
		FeedID: dbFeed.ID,
		UserID: dbUser.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Feed follow is created from feed: %s, user: %s\n", dbFeed.Name, dbUser.Name)
	return nil
}

func handlerfollowing(s *state, cmd command, dbUser database.User) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("incorrect argument for following, need following only")
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.TODO(), dbUser.ID)
	if err != nil {
		return err
	}
	fmt.Printf("User %s is following:\n", s.cfg.CurrentUserName)
	if len(feeds) == 0 {
		fmt.Printf("User %s is not following any feeds", s.cfg.CurrentUserName)
		return nil
	}
	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.String)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, dbUser database.User) error {
	if len(cmd.arguments) != 3 {
		return fmt.Errorf("incorrect argument for following, need unfollow URL")
	}
	err := s.db.DeleteFollowFromURLandUser(context.TODO(), database.DeleteFollowFromURLandUserParams{
		Url:    cmd.arguments[2],
		UserID: dbUser.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Delete following between %s and %s", dbUser.Name, cmd.arguments[2])
	return nil
}
