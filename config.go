package bumble

import "os"

// Config contains the data storage configuration.
type Config struct {
	DatabaseURI string
	ImagesPath  string
}

// GetConfig gets the configuration from the environment,
// using default values if necessary.
func GetConfig() *Config {
	return &Config{
		DatabaseURI: getDatabaseURI(),
		ImagesPath:  getImagesPath(),
	}
}

func getDatabaseURI() string {
	res := os.Getenv("BUMBLE_DB")
	if res != "" {
		return res
	}
	return "mongodb://localhost:27017"
}

func getImagesPath() string {
	res := os.Getenv("BUMBLE_IMAGES")
	if res != "" {
		return res
	}
	return "./photos"
}
