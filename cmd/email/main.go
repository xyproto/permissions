package main

// This is a quick test program in connection with issue #1

import (
	"log"

	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/xyproto/permissions"
	"github.com/xyproto/simpleredis/v2"
)

// Get a string from a list of results at a given position
func getString(bi []interface{}, i int) string {
	return string(bi[i].([]uint8))
}

func main() {
	userstate := permissions.NewUserStateSimple()
	userstate.AddUser("bob", "hunter1", "bob@zombo.com")
	username, err := userstate.HasEmail("bob@zombo.com")
	if err != nil {
		log.Fatalln("Error, the e-mail should exist:", err)
	}
	if username != "bob" {
		log.Fatalln("Error, the e-mail address should belong to bob, but belongs to:", username)
	}
	username, err = userstate.HasEmail("rob@zombo.com")
	if err != permissions.ErrNotFound {
		log.Fatalln("Error, the e-mail should not exist: " + username)
	}

	pool := userstate.Pool()

	// Convert from a simpleredis.ConnectionPool to a redis.Pool
	redisPool := redis.Pool(*pool)
	fmt.Printf("Redis pool: %v (%T)\n", redisPool, redisPool)

	// Get the Redis connection as well
	redisConnection := redisPool.Get()
	fmt.Printf("Redis connection: %v (%T)\n", redisConnection, redisConnection)

	msg, err := redisConnection.Do("PING")
	if err != nil {
		log.Fatalln("could not send PING")
	}
	if msg != "PONG" {
		log.Fatalln("did not get PONG in return")
	}

	users := simpleredis.NewHashMap(pool, "users")
	users.SelectDatabase(0)

	result, err := redis.Values(redisConnection.Do("KEYS", "*"))
	if err != nil {
		log.Fatalln("could not send PING")
	}

	strs := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		strs[i] = getString(result, i)
	}

	fmt.Println(strs)

	email, err := users.Get("bob", "email")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(email)

	username, err = users.FindIDByFieldValue("email", "bob@zombo.com")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(username)
}
