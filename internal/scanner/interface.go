package scanner

// Scanner is the public interface of a scanner object.
type Scanner interface {
	GetObjects() ([]Object, error)
	SetConfig(Config)
	GetConfig() Config
	Scale(Object, int) error
}
