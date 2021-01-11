package config

// Env variables
type Env struct {
	NodeEnv       string `default:"dev" split_words:"true"`
	ServiceName   string `required:"true" split_words:"true"`
	Port          int    `required:"true"`
	ServiceSecret []byte `required:"true" split_words:"true"`

	// mongodb
	MongodbURL  string `required:"true" split_words:"true"`
	MongodbName string `required:"true" split_words:"true"`

	// redis
	RedisHost     string `required:"true" split_words:"true"`
	RedisPort     int    `required:"true" split_words:"true"`
	RedisPassword string `default:"" split_words:"true"`
}
