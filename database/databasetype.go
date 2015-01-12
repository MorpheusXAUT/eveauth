package database

type DatabaseType int

const (
	DatabaseTypeMock DatabaseType = iota - 1
	DatabaseTypeMemory
	DatabaseTypeMySQL
)

func (t DatabaseType) String() string {
	switch t {
	case DatabaseTypeMock:
		return "Mock"
	case DatabaseTypeMemory:
		return "Memory"
	case DatabaseTypeMySQL:
		return "MySQL"
	default:
		return "Unknown"
	}
}
