package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
)

var file = flag.String("f", "", "Execute commands from FILE, then exit.")
var command = flag.String("e", "", "Execute the CQL statement and exit.")
var username = flag.String("u", "cassandra", "Authenticate as user. Default = cassandra.")
var password = flag.String("p", "cassandra", "Authenticate using password. Default = cassandra")
var keyspace = flag.String("k", "", "Use the given keyspace. Equivalent to issuing a USE keyspace command immediately after starting cqlsh.")

func main() {
	//Disable default logger (for gocql)
	//log.SetOutput(ioutil.Discard)
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

	if *command != "" {
		if strings.HasPrefix(strings.ToLower(*command), "select") {
			if rows, err := session.Query(*command).Iter().SliceMap(); err == nil {
				data, err := json.MarshalIndent(rows, "", "  ")
				if err != nil {
					fmt.Println("Failed to format CQL result", err)
					os.Exit(1)
				}
				fmt.Println(string(data))
			} else {
				fmt.Println("Failed to execute CQL command", err)
				os.Exit(1)
			}
			/*
				cols := iter.Columns()
				if len(cols) == 0 {
					fmt.Printf("No result : %v\n", iter.Close())
					return
				}
			*/
			/*
				for {
					// New map each iteration
					row = make(map[string]interface{})
					if !iter.MapScan(row) {
						break
					}
					// Do things with row
					if fullname, ok := row["fullname"]; ok {
						fmt.Printf("Full Name: %s\n", fullname)
					}
				}
			*/
			//fmt.Println("Failed to execute CQL command", err)
			//os.Exit(1)
			//iter.Close()

		} else {
			if err = session.Query(*command).Exec(); err != nil {
				fmt.Println("Failed to execute CQL command", err)
				os.Exit(1)
			}
		}
		fmt.Println("Success !")
	} else if *file != "" {
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
	} else {
		fmt.Println("Nothing todo: -e and -f undefined")
	}
}

func helpMsg() {
	fmt.Println("Usage: gocqlcli [options] [host [port]]")
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(1)
}
