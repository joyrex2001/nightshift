package scanner

// Scale will scale the Object to the given amount of replicas.
func (obj *Object) Scale(replicas int) error {
	scanner, err := New(obj.Type)
	if err != nil {
		return err
	}
	return scanner.Scale(obj, replicas)
}

// SaveState will save the current number of replicas.
func (obj *Object) SaveState() error {
	scanner, err := New(obj.Type)
	if err != nil {
		return err
	}
	return scanner.SaveState(obj)
}

// LoadState will load the number of replicas that has been previously saved.
func (obj *Object) LoadState() error {
	scanner, err := New(obj.Type)
	if err != nil {
		return err
	}
	return scanner.LoadState(obj)
}
