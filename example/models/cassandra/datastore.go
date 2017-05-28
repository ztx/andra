// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=github.com/ztx/andra/example/design
// --out=$(GOPATH)/src/github.com/ztx/andra/example
// --version=v1.1.0-dirty
//
// API "entp": Cassandra DataStore
//
// The content of this file is auto-generated, DO NOT MODIFY

package cassandra

import (
	"github.com/gocql/gocql"
)

//_session holds session
// suffixed with "_" to not to confuse with local names in test files.
var _session *gocql.Session
var (
	ErrSessionAlreadyExists = "Session exits; Only one session is allowed."
)

func init() {
	_session = createSession([]string{"127.0.0.1"}, "storage")
}
func Session() *gocql.Session {
	return _session
}
func createSession(hosts []string, keyspace string) *gocql.Session {
	//if session exits return the same
	if _session != nil {
		return _session
	}
	//else create one
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	_session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Error while creating cassandra session", err)
	}

	return _session
}

func CloseSession() error {
	_session.Close()
	return nil
}
