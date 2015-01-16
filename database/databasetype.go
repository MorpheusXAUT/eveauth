package database

type DatabaseType int

const (
	DatabaseTypeNone DatabaseType = iota
	DatabaseTypeMySQL
)

func (t DatabaseType) String() string {
	switch t {
	case DatabaseTypeMySQL:
		return "MySQL"
	default:
		return "Unknown"
	}
}
