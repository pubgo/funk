package errors

func layerTags(err error) Tags {
	switch e := err.(type) {
	case *ErrWrap:
		return e.Tags
	case *Err:
		return e.Tags
	case Err:
		return e.Tags
	default:
		return nil
	}
}

// GetTags returns tags from the innermost *Err in the chain.
// When no *Err carries tags, it falls back to the outermost layer tags.
func GetTags(err error) Tags {
	if err == nil {
		return nil
	}

	var inner Tags
	for current := err; current != nil; current = Unwrap(current) {
		if e, ok := current.(*Err); ok && len(e.Tags) > 0 {
			inner = e.Tags
		}
	}

	if len(inner) > 0 {
		return cloneTags(inner)
	}

	if tags := layerTags(err); len(tags) > 0 {
		return cloneTags(tags)
	}

	return nil
}

// CollectTags merges tags from every *Err and *ErrWrap layer in the chain.
// Outer layers win when keys conflict.
func CollectTags(err error) Tags {
	if err == nil {
		return nil
	}

	collected := make(Tags)
	for current := err; current != nil; current = Unwrap(current) {
		tags := layerTags(current)
		for key, value := range tags {
			if _, exists := collected[key]; !exists {
				collected[key] = value
			}
		}
	}

	if len(collected) == 0 {
		return nil
	}
	return collected
}
