package errors

func (t *baseErr) Tags() Tags {
	return t.tags
}

func (t *baseErr) AddTag(key string, val any) {
	if t.tags == nil {
		t.tags = make(Tags)
	}
	t.tags[key] = val
}

func (t *baseErr) AddTags(m Tags) {
	if m == nil || len(m) == 0 {
		return
	}

	if t.tags == nil {
		t.tags = m
	}

	for k, v := range m {
		t.tags[k] = v
	}
}
