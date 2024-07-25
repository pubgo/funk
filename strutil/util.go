package strutil

import _ "golang.org/x/text/cases"

func FirstFnNotEmpty(ff ...func() string) string {
	for i := range ff {
		v := ff[i]()
		if v != "" {
			return v
		}
	}
	return ""
}

func FirstNotEmpty(ff ...string) string {
	for i := range ff {
		v := ff[i]
		if v != "" {
			return v
		}
	}
	return ""
}

func GetDefault(names ...string) string {
	name := "default"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return name
}
