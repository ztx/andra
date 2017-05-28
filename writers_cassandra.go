package andra

import (
	"text/template"

	"github.com/goadesign/goa/goagen/codegen"
)

type CassandraWriter struct {
	*codegen.SourceFile
	cassandraTmpl *template.Template
}

type CassandraTmplData struct {
	Cluster  []string
	KeySpace string
}

func NewCassandraWriter(filename string) (*CassandraWriter, error) {
	file, err := codegen.SourceFileFor(filename)
	if err != nil {
		return nil, err
	}
	return &CassandraWriter{SourceFile: file}, nil
}

func (w *CassandraWriter) Execute(data *CassandraTmplData) error {
	fm := make(map[string]interface{})
	fm["wlov"] = writeLOV

	return w.ExecuteTemplate("cassandraTmpl", cassandraDataStoreT, fm, data)
}

const (
	cassandraDataStoreT = `
	//_session holds session
// suffixed with "_" to not to confuse with local names in test files.
var _session *gocql.Session
var (
	ErrSessionAlreadyExists = "Session exits; Only one session is allowed."
)

func init() {
	_session = createSession([]string{ {{range .Cluster}}"{{.}}"{{end}} }, "{{.KeySpace}}")
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
	`
)
