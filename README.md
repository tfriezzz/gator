# gator

gator is a RSS feed aggregator written in **go** with a **postgres** database.


## How to Run Gator
0. You need to have **go**, **postgres** and **goose** installed to run the program.

1. Install gator with `go install github.com/tfriezzz/gator@latest`.

2. Create a new postgres database, preferrably named gator.

3. Create a config file in your home directory, '~/gatorconfig.json', with the following content :

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
}
```

(change the connection string to the appropriate url)


4. Run the database migration with: `goose postgres <connection_string> up`. “Replace `<connection_string>` with the same URL you put in `db_url`.”

5. Run Gator in your terminal with `gator`

Note: During development you can run the project with `go run .`, but once it’s installed with `go install` you should just use the compiled binary: `gator`.


## How to Use Gator

- `register <username>` to register a user
- `login <username>` to login a user
- `users` to list all the users
- `reset` the user list
- `addfeed <name> <url>` to add a rss feed
- `feeds` lists the added feeds
- `follow <url>` to follow another user's feed
- `unfollow <url>` unfollow the feed
- `following` lists the names of the user's followed feeds
- `agg <interval>` starts aggregating posts from the added feeds (run in a separate terminal or as a daemon). Ex: `agg 30m20s`
- `browse <limit>` browses the users feeds
