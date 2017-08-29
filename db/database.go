package db

// A KeyRing contains credentials for connecting to a database.
type KeyRing struct {
	Database string
	Host string
	Port uint
	User string
	Password string
}
