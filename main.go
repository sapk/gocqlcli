package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/gocql/gocql"
)

var file = flag.String("f", "file_name", "Execute commands from FILE, then exit.")
var username = flag.String("u", "cassandra", "Authenticate as user. Default = cassandra.")
var password = flag.String("p", "cassandra", "Authenticate using password. Default = cassandra")
var keyspace = flag.String("k", "", "Use the given keyspace. Equivalent to issuing a USE keyspace command immediately after starting cqlsh.")

func main() {
	//Disable default logger (for gocql)
	log.SetOutput(ioutil.Discard)
	flag.Parse()
	args := flag.Args()
	if len(args) > 2 {
		helpMsg()
	}
	server := "localhost"
	if len(args) > 0 {
		server = args[0]
	}
	port := "9042"
	if len(args) == 2 {
		port = args[1]
	}
	cluster := gocql.NewCluster(server)
	if *keyspace != "" {
		cluster.Keyspace = *keyspace
	}
	if s, err := strconv.ParseInt(port, 10, 64); err == nil {
		cluster.Port = int(s)
	} else {
		fmt.Printf("Invalid port number: %s\n", port)
		helpMsg()
	}

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: *username,
		Password: *password,
	}

	fmt.Printf("Connecting to %s:%s ...\n", server, port)
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		fmt.Println("No connection to cassandra cluster", err)
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println("Failed to read CQL script file", err)
		os.Exit(1)
	}

	if err = session.Query(string(data)).Exec(); err != nil {
		fmt.Println("Failed to execute CQL script file", err)
		os.Exit(1)
	}
	fmt.Println("Success !")
}

func helpMsg() {
	fmt.Println("Usage: gocqlcli [options] [host [port]]")
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}
