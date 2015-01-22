package database

type Type int

const (
	TypeNone Type = iota
	TypeMySQL
)

func (t Type) String() string {
	switch t {
	case TypeMySQL:
		return "MySQL"
	default:
		return "Unknown"
	}
}
