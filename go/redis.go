package openapi

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

// Todo: there might be a better way instead of using a global var
// https://stackoverflow.com/questions/33646948/go-using-mux-router-how-to-pass-my-db-to-my-handlers
var Pool *redis.Pool

// https://medium.com/@gilcrest_65433/basic-redis-examples-with-go-a3348a12878e
func RedisConnect() {
	log.Printf("Redis->connecting...")
	log.Printf("Redis->requesting new pool...")
	// newPool returns a pointer to a redis.Pool
	Pool = redisNewPool()
	redisPing(Pool.Get()) // test connection via ping
}

func redisNewPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80, //80
		// max number of connections
		MaxActive: 12000, //12000
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			log.Printf("Redis->requesting connection\n")
			//Todo: make IP configureable
			c, err := redis.Dial("tcp",
				"35.202.143.31:6379",
				redis.DialConnectTimeout(3*time.Second))
			if err != nil {
				log.Printf("Redis->cannot connect :-(\n")
				panic(err.Error())
			}

			//Todo: don't set PW in code ;-)
			response, err := c.Do("AUTH", "MKX9xoTPT8Ca") //redis1
			//response, err := c.Do("AUTH", "CX9EN8as6UBS")	//redis2
			log.Printf("Redis->connection: %s\n", response)
			return c, err
		},
	}
}

// https://medium.com/@gilcrest_65433/basic-redis-examples-with-go-a3348a12878e
// ping tests connectivity for redis (PONG should be returned)
func redisPing(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	log.Printf("Redis->Ping: %s\n", s)
	// Output: PONG

	return nil
}

func redisQueryEquipmentStatus(pei string) (EquipmentStatus, error) {
	// get a connection from the pool (redis.Conn)
	var conn redis.Conn
	conn = Pool.Get()
	// use defer to close the connection when the function completes
	defer conn.Close()
	val, err := redisGet(conn, pei)
	if err != nil {
		return "", err
	}
	return val, err

}

// get executes the redis GET command
func redisGet(c redis.Conn, pei string) (EquipmentStatus, error) {
	//pei := "imei-123456789012345"
	//s, err := redis.String(c.Do("GET", pei))
	s, err := redis.String(c.Do("GET", pei))
	if err != nil {
		log.Printf("Redis->error get EquipmentStatus: %s\n", err)
		return "", err
	}

	return EquipmentStatus(s), nil
}
