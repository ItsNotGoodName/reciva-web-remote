package state

func (c Changed) Is(check Changed) bool {
	return c&check == check
}

func (c Changed) Merge(merge Changed) Changed {
	return c | merge
}
