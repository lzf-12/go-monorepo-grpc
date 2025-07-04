package env

type Environment int

const (
	Dev Environment = iota
	Staging
	Production
)

func (e Environment) String() string {
	switch e {
	case Dev:
		return "dev"
	case Staging:
		return "staging"
	case Production:
		return "production"
	default:
		return "unknown"
	}
}
