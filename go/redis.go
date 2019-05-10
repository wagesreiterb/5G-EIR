package openapi

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

var Pool *redis.Pool
var Conn redis.Conn

/*
// **************************************************************************
// Todo:
// the redis instance shouldn't be created for every query newly
// otherwise a query takes quite long ~300ms
// https://stackoverflow.com/questions/43695832/pass-a-reference-to-a-redis-instance-to-a-gorilla-mux-handler
// **************************************************************************
func redisQueryEquipmentStatus(pei string) (EquipmentStatus, error) {
	// newPool returns a pointer to a redis.Pool
	pool := RedisNewPool()
	// get a connection from the pool (redis.Conn)
	conn := pool.Get()
	// use defer to close the connection when the function completes
	defer conn.Close()

	// set demonstrates the redis GET command
	val, err := redisGet(conn, pei)
	if err != nil {
		return "", err
	}
	fmt.Println("pei2:", pei)
	return val, err

}
*/

func RedisConnect() {
	log.Printf("Redis->Connecting...")
	// newPool returns a pointer to a redis.Pool
	Pool = RedisNewPool()
	// get a connection from the pool (redis.Conn)
	Conn = Pool.Get()
	// use defer to close the connection when the function completes
	//defer Conn.Close()
}

func redisQueryEquipmentStatus(pei string) (EquipmentStatus, error) {
	// set demonstrates the redis GET command
	val, err := redisGet(Conn, pei)
	if err != nil {
		return "", err
	}
	return val, err

}

// get executes the redis GET command
// func redisGet(c redis.Conn, pei string) (EquipmentStatus, error) {
func redisGet(c redis.Conn, pei string) (EquipmentStatus, error) {
	//pei := "imei-123456789012345"
	s, err := redis.String(Conn.Do("GET", pei))
	if err != nil {
		return "", err
	}

	return EquipmentStatus(s), nil
}

func RedisNewPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "35.238.180.253:6379")
			if err != nil {
				panic(err.Error())
			}

			response, err := c.Do("AUTH", "MKX9xoTPT8Ca")
			//fmt.Fprintf(w, "Connected: %s\n", response)
			log.Printf("Redis->Connected: %s\n", response)
			return c, err
		},
	}
}
