package utils

import "fmt"

func Must[T any](v T, err error) T {
	if err != nil { panic(err) }
	return v
}

func DSN(host string, port int, user, pass, db, sslmode string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, pass, host, port, db, sslmode)
}
