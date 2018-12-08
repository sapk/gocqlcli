# gocqlcli
A simple wrapper around https://github.com/gocql/gocql to execute script file (try to keep compat args with cqlsh)

```
Usage: gocqlcli [options] [host [port]]
Options:
  -e string
    	Execute the CQL statement and exit.
  -f string
    	Execute commands from FILE, then exit.
  -k string
    	Use the given keyspace. Equivalent to issuing a USE keyspace command immediately after starting cqlsh.
  -p string
    	Authenticate using password. Default = cassandra (default "cassandra")
  -u string
    	Authenticate as user. Default = cassandra. (default "cassandra")
```
