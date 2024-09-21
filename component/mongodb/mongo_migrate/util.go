package mongo_migrate

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pubgo/funk/errors"
)

func extractVersionDescription(name string) (uint64, string, error) {
	base := filepath.Base(name)

	if ext := filepath.Ext(base); ext != ".go" {
		return 0, "", errors.Format("can not extract version from %q", base)
	}

	idx := strings.IndexByte(base, '_')
	if idx == -1 {
		return 0, "", errors.Format("can not extract version from %q", base)
	}

	version, err := strconv.ParseUint(base[:idx], 10, 64)
	if err != nil {
		return 0, "", err
	}

	description := base[idx+1 : len(base)-len(".go")]

	return version, description, nil
}
