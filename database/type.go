package database

// Type represents the type of database backend to use
type Type int

const (
	// TypeNone represents a non-persistent database backend, only storing values in memory (not implemented)
	TypeNone Type = iota
	// TypeMySQL represents a persistent MySQL database backend
	TypeMySQL
)

// String returns a easily readable string representations of the given Type
func (t Type) String() string {
	switch t {
	case TypeMySQL:
		return "MySQL"
	default:
		return "Unknown"
	}
}
