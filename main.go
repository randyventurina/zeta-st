package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	utils "utils"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

//Node is a node within zetanet
type Node struct {
	Name    string
	Host    string
	Port    string
	Type    string
	Country string
	Path    string
}

//Desc is the description structure of self-describing hash
type Desc struct {
	fp string //fn file_path
	ha string //hs hashing_algorithm
	le string //le content_hash_length
	ep string //ep endpoint
}

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	register()
	command()
}

func command() {
	//add subcommand
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addFile := addCmd.String("file", "", "Base64 Encoding of file ")
	addPush := addCmd.Bool("push", true, "Pushes contents to zetanet ")

	//remove subcommand
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	removeFilePath := removeCmd.String("filePath", "", "Path of the file to be removed")

	//update file subcommand
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateFile := updateCmd.String("file", "", "Base64 Encoding of file ")
	updateDest := updateCmd.String("dest", "", "File path to be updated")

	//get file subcommand
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getFilePath := getCmd.String("filePath", "", "The path or url of the file to be retrieved")

	//hash file subcommand
	hashCmd := flag.NewFlagSet("hash", flag.ExitOnError)
	hashFilePath := hashCmd.String("filePath", "", "The path or url of the file to be retrieved")

	if len(os.Args) < 2 {
		fmt.Println("expected 'add|update|remove|get' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		hash := add(*addFile, *addPush)
		fmt.Println(hash)

	case "remove": // not applicable anymore
		removeCmd.Parse(os.Args[2:])
		fmt.Println("Remove Command")
		fmt.Println("   filePath:", *removeFilePath)
	case "update": // not applicable anymore
		updateCmd.Parse(os.Args[2:])
		fmt.Println("Update Command")
		fmt.Println("   file:", *updateFile)
		fmt.Println("   dest:", *updateDest)
	case "get":
		getCmd.Parse(os.Args[2:])
		fmt.Println("Get Command")
		fmt.Println("   filePath:", *getFilePath)
	case "hash":
		hashCmd.Parse(os.Args[2:])
		fmt.Println("Hash Command")
		fmt.Println("	filePath: ", *hashFilePath)
		if hash, _ := utils.Hash("md5", *hashFilePath); hash != "" {
			fmt.Println("	Content Hash: " + hash)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
	os.Exit(0)
}

//saves the content-address to local leveldb
//copy the file to universe folder
func add(file string, push bool) string {
	

	addToUniverse(file)	

	// save hash to local leveldb
	if push {
		hash, err := saveHashLocally(file)
		
		if err == nil {
			return hash
		}

		saveHashGlobally()
	}

	return "content hashing: failed"
}

func saveHashGlobally() (string, error){
	return "", nil
}

//save a content hash entry to leveldb
func saveHashLocally(file string) (string, error) {
	hash, _ := utils.Hash("md5", file)

	db, err := leveldb.OpenFile("db.contents", nil)
	desc, err := describeInBytes(file)

	db.Put([]byte(hash), desc, nil)
	db.Close()
	
	if err == nil {
		return hash, err 
	}

	return hash, err
}

func addToUniverse(file string) bool { 
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fileName := filepath.Base(file)
	utils.Copy(file, filepath.Join(exPath, "/universe", fileName) )
	return true
}

//hashes the self-describing part of content-address
func hashDescription(desc *Desc) string{
	b, _:= json.Marshal(desc)
	return utils.HashContent(b)
}

//creates description of the file being hashed
func describeInObject(file string) (*Desc, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
	  // Could not obtain stat, handle error
	}

	var stConfig Config
	stConfig.Init("./config/st." + utils.LoadEnv() + ".yml")

	endpoint := stConfig.Host + ":" + stConfig.Port

	desc := Desc{fp:file, ha:"md5", le:string(fi.Size()), ep:endpoint }

	return &desc, nil
}

func describeInBytes(file string) ([]byte, error){
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
	  // Could not obtain stat, handle error
	}

	var stConfig Config
	stConfig.Init("./config/st." + utils.LoadEnv() + ".yml")

	endpoint := stConfig.Host + ":" + stConfig.Port

	desc := Desc{fp:file, ha:"md5", le:string(fi.Size()), ep:endpoint }

	b, _:= json.Marshal(desc)

	return b, nil
}

//register to discovery node. registration is required, this is to make sure all nodes are discovered within the network
func register() {

	// registration to disovery node must run once
	o := opt.Options{ErrorIfExist: true}
	db, err := leveldb.OpenFile("db.nodes", &o)

	// run once
	if err == nil {
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
				saveNode(db, conn, message)
			}
		} else {
			fmt.Println(err)
		}
	}
}

//save initial list of nodes upon successful connection
func saveNode(db *leveldb.DB, conn net.Conn, data []byte) bool {
	var node Node

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

	return true
}
