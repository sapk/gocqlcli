package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

var (
	//Version version of app set by build flag
	Version string
	//Branch git branch of app set by build flag
	Branch string
	//Commit git commit of app set by build flag
	Commit string
	//BuildTime build time of app set by build flag
	BuildTime string
)

var displayVersion = flag.Bool("version", false, "gocqlcli version.")
var file = flag.String("f", "", "Execute commands from FILE, then exit.")
var command = flag.String("e", "", "Execute the CQL statement and exit.")
var username = flag.String("u", "cassandra", "Authenticate as user. Default = cassandra.")
var password = flag.String("p", "cassandra", "Authenticate using password. Default = cassandra")
var keyspace = flag.String("k", "", "Use the given keyspace. Equivalent to issuing a USE keyspace command immediately after starting cqlsh.")

func main() {
	//Disable default logger (for gocql)
	//log.SetOutput(ioutil.Discard)
	flag.Parse()
	if *displayVersion {
		displayVersionMsg()
		return
	}
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
	cluster.Timeout = 30 * time.Second
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

	if *command != "" {
		executeCQL(cluster, *command)
	} else if *file != "" {
		data, err := ioutil.ReadFile(*file)
		if err != nil {
			fmt.Println("Failed to read CQL script file", err)
			os.Exit(1)
		}
		commands := strings.Split(string(data), ";")
		for _, c := range commands {
			c = strings.TrimSpace(c)
			if c != "" {
				executeCQL(cluster, c+";")
			}
		}
	} else {
		fmt.Println("Nothing todo: -e and -f undefined")
	}
}

func executeCQL(cluster *gocql.ClusterConfig, cmd string) {
	//TODO ignore ErrTimeoutNoResponse
	//cmd = strings.Replace(cmd, "\n", "", -1) //TODO better
	//cmd = strings.Replace(cmd, "\r", "", -1) //TODO better
	fmt.Println("Executing CQL command", cmd)
	if strings.HasPrefix(strings.ToLower(cmd), "use") { //Change keyspace
		*keyspace = strings.ToLower(strings.TrimSpace(strings.Trim(cmd[3:], ";")))
		//time.Sleep(30 * time.Second) //TODO better to wait for keyspace if is in creation
		fmt.Println("Success !")
		return
	} else {
		if *keyspace != "" {
			cluster.Keyspace = *keyspace
		}
		session, err := cluster.CreateSession()
		defer session.Close()
		if err != nil {
			fmt.Println("No connection to cassandra cluster", err)
			os.Exit(1)
		}
		if strings.HasPrefix(strings.ToLower(cmd), "select") {
			if rows, err := session.Query(cmd).Iter().SliceMap(); err == nil {
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
		} else {
			if err := session.Query(cmd).Exec(); err != nil {
				fmt.Println("Failed to execute CQL command", err)
				os.Exit(1)
			}
		}
	}
	fmt.Println("Success !")
}

func displayVersionMsg() {
	fmt.Printf("\nVersion: %s - Branch: %s - Commit: %s - BuildTime: %s\n\n", Version, Branch, Commit, BuildTime)
}

func helpMsg() {
	fmt.Println("Usage: gocqlcli [options] [host [port]]")
	fmt.Println("Options:")
	flag.PrintDefaults()
	displayVersionMsg()
	os.Exit(1)
}
