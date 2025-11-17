// Package handleruser handles user specific actions
package handleruser

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	// "github.com/tfriezzz/gator/internal/commands"
	"github.com/google/uuid"
	"github.com/tfriezzz/gator/commands"
	"github.com/tfriezzz/gator/internal/config"
	"github.com/tfriezzz/gator/internal/database"
	"github.com/tfriezzz/gator/rss"
)

func requireOneArg(cmd commands.Command, usage string) {
	if len(cmd.Args) < 1 {
		fmt.Fprintln(os.Stderr, "usage:", usage)
		os.Exit(1)
	}
}

// func getCurrentUser(s *config.State) database.User {
// 	userName := s.Config.CurrentUserName
// 	currentUser, err := s.DB.GetUser(context.Background(), userName)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	return currentUser
// }

func MiddlewareLoggedIn(handler func(s *config.State, cmd commands.Command, user database.User) error) func(*config.State, commands.Command) error {
	newHandler := func(s *config.State, cmd commands.Command) error {
		currentUser, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			fmt.Println(err)
		}

		handler(s, cmd, currentUser)

		return nil
	}

	return newHandler
}

func scrapeFeeds(s *config.State, cmd commands.Command) error {
	feed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	feedFetchedParams := database.MarkFeedFetchedParams{
		feed.ID, feed.LastFetchedAt, feed.UpdatedAt,
	}
	s.DB.MarkFeedFetched(context.Background(), feedFetchedParams)

	fetchedFeed, err := rss.FetchFeed(context.Background(), feed.Url)

	for _, RSSItem := range fetchedFeed.Channel.Item {
		fmt.Println(RSSItem.Title)
	}
	return nil
}

func HandlerLogin(s *config.State, cmd commands.Command) error {
	requireOneArg(cmd, "gator login <username>")

	userName := cmd.Args[0]
	_, err := s.DB.GetUser(context.Background(), userName)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Fprintln(os.Stderr, "user not found")
		os.Exit(1)
	}
	err = s.Config.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Printf("the user has been set to %s", userName)

	return nil
}

func HandlerRegister(s *config.State, cmd commands.Command) error {
	requireOneArg(cmd, "gator register <username>")
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the register handler expects a single argument")
	}

	userID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	userName := cmd.Args[0]

	userParams := database.CreateUserParams{
		userID, createdAt,
		updatedAt, userName,
	}

	_, err := s.DB.GetUser(context.Background(), userName)
	if err == nil {
		fmt.Fprintln(os.Stderr, "user already exists")
		os.Exit(1)
	}

	user, err := s.DB.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	if err = s.Config.SetUser(user.Name); err != nil {
		os.Exit(1)
		return err
	}

	fmt.Printf("user %s was created: %+v", user.Name, user)
	return nil
}

func HandlerReset(s *config.State, cmd commands.Command) error {
	err := s.DB.DeleteUser(context.Background())
	if err != nil {
		return err
	}
	fmt.Print("user list has been reset")
	return nil
}

func HandlerList(s *config.State, cmd commands.Command) error {
	currentUser := s.Config.CurrentUserName
	list, err := s.DB.ListUser(context.Background())
	if err != nil {
		return err
	}
	for _, u := range list {
		isCurrentUser := ""
		if u.Name == currentUser {
			isCurrentUser = "(current)"
		}
		fmt.Printf("* %v %v\n", u.Name, isCurrentUser)
	}
	return nil
}

func HandlerAgg(s *config.State, cmd commands.Command) error {
	requireOneArg(cmd, "gator agg <time_between_reqs>")
	timeBetweenRequests, err := time.ParseDuration(os.Args[2])
	if err != nil {
		return err
	}
	fmt.Printf("collecting feeds every %v\n", os.Args[2])

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s, cmd)
		if err != nil {
			return err
		}
	}

	// testFeed, _ := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	// fmt.Printf("testFeed: %v", testFeed)

	return nil
}

func HandlerAddFeed(s *config.State, cmd commands.Command, u database.User) error {
	// userName := s.Config.CurrentUserName
	// currentUser, err := s.DB.GetUser(context.Background(), userName)
	// if err != nil {
	// 	return err
	// }

	ID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	name := cmd.Args[0]
	url := cmd.Args[1]

	feedParams := database.CreateFeedParams{
		ID, createdAt, updatedAt, name, url, u.ID,
	}

	feed, err := s.DB.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	fmt.Println(feed)

	feedFollowParams := database.CreateFeedFollowParams{
		feed.ID, feed.CreatedAt, feed.UpdatedAt, u.ID, feed.ID,
	}

	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("feed '%v' followed", feedFollow.FeedName)

	return nil
}

func HandlerFeeds(s *config.State, cmd commands.Command) error {
	feeds, err := s.DB.ListFeed(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, _ := s.DB.GetUserByID(context.Background(), feed.UserID)

		fmt.Printf("%v, %v\n", feed, user.Name)
	}

	return nil
}

func HandlerFollow(s *config.State, cmd commands.Command, u database.User) error {
	url := cmd.Args[0]
	feed, err := s.DB.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	ID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	userID := u.ID
	feedID := feed.ID

	feedFollowParams := database.CreateFeedFollowParams{
		ID, createdAt, updatedAt, userID, feedID,
	}

	feedFollow, err := s.DB.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("feed: %v, user: %v", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func HandlerFollowing(s *config.State, cmd commands.Command, u database.User) error {
	feeds, err := s.DB.GetFeedFollowsForUser(context.Background(), u.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}

	return nil
}

func HandlerUnfollow(s *config.State, cmd commands.Command, u database.User) error {
	url := cmd.Args[0]
	feed, err := s.DB.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	unfollowParams := database.DeleteFeedFollowParams{
		u.ID, feed.ID,
	}
	s.DB.DeleteFeedFollow(context.Background(), unfollowParams)

	return nil
}
