package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	usersCount   = 100
	workersCount = 100
)

var actions = []string{"logged in", "logged out", "created record", "deleted record", "updated account"}

type logItem struct {
	action    string
	timestamp time.Time
}

type User struct {
	id    int
	email string
	logs  []logItem
}

func (u User) getActivityInfo() string {
	output := fmt.Sprintf("UID: %d; Email: %s;\nActivity Log:\n", u.id, u.email)
	for index, item := range u.logs {
		output += fmt.Sprintf("%d. [%s] at %s\n", index, item.action, item.timestamp.Format(time.RFC3339))
	}

	return output
}

func main() {
	wg := sync.WaitGroup{}

	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	generatedUsers := make(chan int, usersCount)
	users := make(chan User, usersCount)

	for i := 0; i < workersCount; i++ {
		go generateUsers(generatedUsers, users)
	}

	worker(usersCount, generatedUsers, &wg)

	for i := 0; i < workersCount; i++ {
		go saveUserInfo(users, &wg)
	}

	wg.Wait()

	fmt.Printf("DONE! Time Elapsed: %.2f seconds\n", time.Since(startTime).Seconds())
}

func worker(count int, generatedUsers chan<- int, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		wg.Add(1)
		generatedUsers <- i
	}
}

func saveUserInfo(users chan User, wg *sync.WaitGroup) {
	for user := range users {
		fmt.Printf("WRITING FILE FOR UID %d\n", user.id)

		filename := fmt.Sprintf("users/uid%d.txt", user.id)
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		file.WriteString(user.getActivityInfo())
		time.Sleep(time.Second)

		wg.Done()
	}
}

func generateUsers(generatedUsers <-chan int, users chan<- User) {
	for i := range generatedUsers {
		users <- User{
			id:    i + 1,
			email: fmt.Sprintf("user%d@company.com", i+1),
			logs:  generateLogs(rand.Intn(1000)),
		}
		fmt.Printf("generated user %d\n", i+1)
		time.Sleep(time.Millisecond * 100)
	}
	close(users)
}

func generateLogs(count int) []logItem {
	logs := make([]logItem, count)

	for i := 0; i < count; i++ {
		logs[i] = logItem{
			action:    actions[rand.Intn(len(actions)-1)],
			timestamp: time.Now(),
		}
	}

	return logs
}
