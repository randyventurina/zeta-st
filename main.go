package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	utils "utils"
)

//Node is a node within zetanet
type Node struct {
	Name    string
	Host    string
	Port    string
	Type    string
	Country string
	Path 	string
}

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	var config Config
	config.Init("./config/dn." + utils.LoadEnv() + ".yml")
	if conn, err := net.Dial(config.Type, config.Host+":"+config.Port); err == nil {
		fmt.Println("Connected to discovery node " + config.Name + " via " + config.Type + " endpoint " + config.Host + ":" + config.Port)

		var stConfig Config
		stConfig.Init("./config/st." + utils.LoadEnv() + ".yml")
		data, _ := json.Marshal(stConfig)
		// send to socket
		conn.Write(append(data, '\n'))

		// listen for reply
		reader := bufio.NewReader(conn)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			message := scanner.Bytes()
			fmt.Println("Node: " + string(message))
			save(conn, message)
		}
	}else {
		fmt.Println(err)
	}
}

func save(conn net.Conn, data []byte) bool {
	var node Node

	if db, err := NewDb(); err == nil {
		defer db.Close()

		if err := json.Unmarshal(data, &node); err == nil {

			//add bytes converted node information to leveldb
			db.Put([]byte(node.Host+":"+node.Port), data, nil)

			//log newly added node
			if data, err := db.Get([]byte(node.Host+":"+node.Port), nil); err == nil {
				fmt.Println("Added node: " + string(data))
			} else {
				fmt.Println("Get:" + err.Error())
			}
		} else {
			fmt.Println("Get:" + err.Error())
		}
	} else {
		fmt.Println("OpenFile:" + err.Error())
	}
	return true
}
