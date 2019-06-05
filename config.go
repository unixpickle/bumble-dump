package bumble

import "os"

// Config contains the data storage configuration.
type Config struct {
	DatabaseURI string
	PhotosPath  string
}

// GetConfig gets the configuration from the environment,
// using default values if necessary.
func GetConfig() *Config {
	return &Config{
		DatabaseURI: getDatabaseURI(),
		PhotosPath:  getPhotosPath(),
	}
}

func getDatabaseURI() string {
	res := os.Getenv("BUMBLE_DB")
	if res != "" {
		return res
	}
	return "mongodb://localhost:27017"
}

func getPhotosPath() string {
	res := os.Getenv("BUMBLE_PHOTOS")
	if res != "" {
		return res
	}
	return "./photos"
}
