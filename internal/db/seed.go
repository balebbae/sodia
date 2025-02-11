package db

import (
	"context"
	"fmt"
	"log"

	"math/rand"

	"github.com/balebbae/sodia/internal/store"
)

var usernames = []string{
    "alice", "bob", "charlie", "dave", "eve", "frank", "grace", "hank", "ivy", "jack",
    "karen", "leo", "mona", "nate", "olivia", "paul", "quincy", "rachel", "steve", "tina",
    "ulysses", "victor", "wanda", "xander", "yasmine", "zane", "aaron", "beth", "carter", "diana",
    "elijah", "fiona", "george", "hannah", "isaac", "jasmine", "kevin", "luna", "miles", "nora",
    "oscar", "penny", "quentin", "ruby", "simon", "tracy", "ursula", "vincent", "wendy", "xavier",
}

var titles = []string{
	"Slice Basics",
    "Growing Slices",
    "Slicing Arrays",
    "Slice Capacity",
    "Slice vs Array",
    "Slice Initialization",
    "Appending Elements",
    "Copying Slices",
    "Reslicing Tricks",
    "Nil vs Empty Slice",
    "Slice Memory Layout",
    "Slice Performance",
    "Subslice Views",
    "Slice Iteration",
    "Slice of Structs",
    "Multi-Dimensional Slices",
    "Sorting Slices",
    "Removing Elements",
    "Passing Slices",
    "Slice Best Practices",
}

var content = []string{
	"In this post, we explore Go routines and how they enable concurrent programming.",
	"Web development is evolving rapidly. Here's what to expect in the coming years.",
	"Artificial Intelligence is transforming industries. What’s next in 2025?",
	"Learn how to build scalable and maintainable APIs using the Go programming language.",
	"TypeScript is gaining popularity among developers. Find out why!",
	"Cyber threats are increasing. Here’s what you need to know to stay safe.",
	"Understanding the core concepts of Machine Learning with real-world applications.",
	"SQL performance is crucial. Learn techniques to optimize your database queries.",
	"GraphQL or REST? We break down the differences and use cases for both.",
	"Blockchain is more than just cryptocurrency. Here’s how it’s changing industries.",
	"The cloud is revolutionizing computing. Learn about the latest trends in cloud technology.",
	"Functional programming can improve your code quality. Here’s how it works.",
	"UI/UX design principles have evolved. Discover the latest best practices.",
	"Writing clean and maintainable code is essential for software projects.",
	"Docker simplifies development and deployment. Learn how to get started.",
	"Python automation can save time and effort. Here’s how to automate tasks.",
	"Kubernetes is the future of app deployment. Learn the basics to get started.",
	"React Hooks have changed the way we write React applications. Let’s explore them.",
	"Managing state in large applications can be challenging. Here are some solutions.",
	"Why data structures and algorithms are fundamental to software engineering.",
}

var tags = []string{
	"Technology", "Programming", "GoLang", "Web Development", "AI", 
	"Machine Learning", "Cybersecurity", "Cloud Computing", "Blockchain", "DevOps", 
	"Databases", "Software Engineering", "Frontend", "Backend", "Mobile Development", 
	"Data Science", "Open Source", "Networking", "Startup", "Automation",
}

var comments = []string{
	"Great post! Really insightful.",
	"I learned a lot from this, thanks!",
	"Can you explain more about this topic?",
	"This was exactly what I was looking for.",
	"Amazing content, keep it up!",
	"I disagree with some points, but overall a good read.",
	"This helped me understand the concept better.",
	"Looking forward to more posts like this.",
	"Very well written and easy to follow.",
	"This clarified so many things for me.",
	"How would this work in a real-world scenario?",
	"I appreciate the detailed explanation.",
	"Could you provide some code examples?",
	"This article is a lifesaver, thanks!",
	"I shared this with my colleagues, very useful!",
	"Do you have any recommended resources on this?",
	"Well explained! I finally get it now.",
	"Nice breakdown of the topic.",
	"This post should be more widely shared!",
	"Great insights! Learned something new today.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		err := store.Users.Create(ctx, user)
		if err != nil {
			log.Println("Errpor creating users:", err)
			return
		}
	}
	
	posts := generatePosts(200, users)
	for _, post := range posts {
		err := store.Posts.Create(ctx, post)
		if err != nil {
			log.Println("Error creating post:", err)
			return 
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		err := store.Comments.Create(ctx, comment)
		if err != nil {
			log.Println("Error creating comments:", err)
			return 
		}
	}

	log.Println("Seeding complete")
}


func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email: usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "1234567",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID: user.ID,
			Title: titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(content))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID: posts[rand.Intn(len(posts))].ID,
			UserID: users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}

