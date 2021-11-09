package config

type User struct {
	Email string `yaml:"email"`
	// UserId is generated from the email stored in the config file we're not saving UserId back to config file
	// so nobody would be tempted to change that in the config file and get inconsistency in between email and userid
	Id string `yaml:"-"`
}

type Format struct {
	Format string `yaml:"format"`
}
type Config struct {
	User   User `yaml:"user"`
	Format Format
}
