package database

type RedisConnectionInfo struct {
	Size     int
	Network  string
	Port     int
	Password string
	Key      string
}
