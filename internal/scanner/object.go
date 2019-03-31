package scanner

// getScanner will lazy load the appropriate scanner object for this resource.
func (obj *Object) getScanner() (Scanner, error) {
	var err error
	if obj.scanner == nil {
		obj.scanner, err = New(obj.Type)
		if err != nil {
			return nil, err
		}
	}
	return obj.scanner, nil
}

// Scale will scale the Object to the given amount of replicas.
func (obj *Object) Scale(replicas int) error {
	scanner, err := obj.getScanner()
	if err != nil {
		return err
	}
	return scanner.Scale(obj, replicas)
}

// SaveState will save the current number of replicas.
func (obj *Object) SaveState() error {
	scanner, err := obj.getScanner()
	if err != nil {
		return err
	}
	return scanner.SaveState(obj)
}

// LoadState will load the number of replicas that has been previously saved.
func (obj *Object) LoadState() error {
	scanner, err := obj.getScanner()
	if err != nil {
		return err
	}
	return scanner.LoadState(obj)
}
