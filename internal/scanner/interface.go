package scanner

// Scanner is the public interface of a scanner object.
type Scanner interface {
	SetConfig(Config)
	GetConfig() Config
	GetObjects() ([]*Object, error)
	SaveState(*Object) error
	LoadState(*Object) error
	Scale(*Object, int) error
}
