package openapi

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// Todo: there might be a better way instead of using a global var
// https://stackoverflow.com/questions/33646948/go-using-mux-router-how-to-pass-my-db-to-my-handlers
var Pool *redis.Pool

// https://medium.com/@gilcrest_65433/basic-redis-examples-with-go-a3348a12878e
func RedisConnect() {
	WriteLog("redis", "connecting...")
	WriteLog("redis", "requesting new pool...")
	Pool = redisNewPool() // returns a pointer to a redis.Pool
	redisPing(Pool.Get()) // test connection via ping
}

func redisNewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,    // Maximum number of idle connections in the pool.
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			WriteLog("redis", "requesting connection")
			//Todo: make IP configureable
			c, err := redis.Dial("tcp",
				//"five-g-eir-redis:6379",	//kubernets
				"five-g-eir-redis:6379",
				//"35.239.35.61:6379",	// VM
				redis.DialConnectTimeout(3*time.Second))
			if err != nil {
				WriteLog("redis", "cannot connect :-(")
				panic(err.Error())
			}

			//Todo: don't set PW in code ;-)
			//response, err := c.Do("AUTH", "MKX9xoTPT8Ca") //redis1
			//response, err := c.Do("AUTH", "CX9EN8as6UBS")	//redis2
			//authResponse, err := redis.String(response, err)
			if err != nil {
				panic(err.Error())
			}
			//WriteLog("redis", "authentication response response: " + authResponse)
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
	pingResponse, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	//log.Printf("Redis->Ping: %s\n", ping_response)
	WriteLog("redis", "ping response: "+pingResponse) // Output if connection works: PONG

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
		WriteLog("redis", "[ERROR] get EquipmentStatus: "+err.Error())
		return "", err
	}

	return EquipmentStatus(s), nil
}
