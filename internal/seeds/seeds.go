package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/sirUnchained/udemy-backend-course/internal/store"
)

var usernames = []string{
	"bluecat42", "greendog88", "redbird17", "fastcar305", "slowhouse99", "bigtree64", "smallbook23", "happysun71",
	"saddog456", "coolstar12", "hotmoon89", "newriver34", "oldmountain67", "smartfood95", "dumbgame28",
	"bluebird156", "greenhouse49", "redcat783", "fastdog21", "slowstar634", "bighouse82", "smallriver37",
	"happybook94", "sadmoon16", "coolcar508", "hotfood73", "newstar29", "olddog661", "smartcat45",
	"dumbhouse98", "bluegame32", "greenriver57", "redsun84", "fastbook19", "slowmountain402", "bigcat76",
	"smallstar63", "happydog31", "sadhouse95", "coolfood48", "hotriver22", "newbook79", "oldstar365",
	"smartdog14", "dumbcat87", "bluehouse51", "greenmountain93", "redgame68", "fastsun35", "slowfood102",
	"bigstar47", "smallcat89", "happyhouse24",
}
var adjectives = []string{"Adventurous", "Creative", "Curious", "Friendly", "Helpful", "Honest", "Kind", "Loyal", "Patient", "Reliable"}
var hobbies = []string{"loves hiking", "enjoys reading", "plays guitar", "codes for fun", "travels often", "cooks meals", "photographs nature", "paints landscapes", "writes stories", "studies history"}
var qualities = []string{"always learning", "seeking challenges", "making friends", "exploring new places", "helping others", "sharing knowledge", "building things", "solving problems", "creating art", "teaching skills"}

func Seed(st store.Storage, debugMode bool, db *sql.DB) {
	if !debugMode {
		log.Println("we are not in debug mode, so ignore seeds")
		return
	}
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := st.Users.Create(ctx, tx, user); err != nil {
			tx.Rollback()
			log.Println("Error: ", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(100, users)
	for _, post := range posts {
		if err := st.Posts.Create(ctx, post); err != nil {
			log.Println("Error: ", err)
			return
		}
	}

	comments := generateComments(100, posts, users)
	for _, comment := range comments {
		if err := st.Comments.Create(ctx, comment); err != nil {
			log.Println("Error: ", err)
			return
		}
	}

}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			UserName: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
		}
		users[i].Password.Set("12341234")
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	// Randomly create a content
	for i := 0; i < num; i++ {
		content := adjectives[rand.Intn(len(adjectives))] + " person who " +
			hobbies[rand.Intn(len(hobbies))] + " and " +
			qualities[rand.Intn(len(qualities))] + "."
		title := qualities[rand.Intn(len(adjectives))] + " and " + hobbies[rand.Intn(len(hobbies))]
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   title,
			Content: content,
			User:    *user,
			Version: 0,
			Tags:    []string{"t1", "21", "t3"},
		}
	}

	return posts
}

func generateComments(num int, posts []*store.Post, users []*store.User) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		content := adjectives[rand.Intn(len(adjectives))] + " person who " +
			hobbies[rand.Intn(len(hobbies))] + " and " +
			qualities[rand.Intn(len(qualities))] + "."
		user := users[rand.Intn(len(users))]
		comments[i] = &store.Comment{
			UserID:  user.ID,
			PostID:  posts[i].ID,
			Content: content,
			User:    *user,
		}
	}

	return comments
}
