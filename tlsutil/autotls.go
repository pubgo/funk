package tlsutil

import (
	"crypto/tls"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/acme/autocert"
)

func NewTlsConfig(domains ...string) *tls.Config {
	m := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}
	if len(domains) > 0 {
		m.HostPolicy = autocert.HostWhitelist(domains...)
	}

	dir := cacheDir()
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		log.Printf("warning: autocert.NewListener not using a cache: %v", err)
	} else {
		m.Cache = autocert.DirCache(dir)
	}
	return m.TLSConfig()
}

func getCacheDir() (autocert.DirCache, error) {
	dir := cacheDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", errors.New("warning: autocert.NewListener not using a cache: " + err.Error())
	}
	return autocert.DirCache(dir), nil
}

func cacheDir() string {
	const base = "golang-autocert"
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir(), "Library", "Caches", base)
	case "windows":
		for _, ev := range []string{"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				return filepath.Join(v, base)
			}
		}
		// Worst case:
		return filepath.Join(homeDir(), base)
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(homeDir(), ".cache", base)
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "/"
}
