package config

import (
	"log"
	"os"
)

type Config struct {
	Addr        string
	DatabaseURL string
	JWTSecret   string
}

func FromEnv() Config {
	c := Config{
		Addr:        getenv("ADDR", ":8080"),
		DatabaseURL: must("DATABASE_URL"),
		JWTSecret:   must("JWT_SECRET"),
	}
	return c
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" { log.Fatalf("missing required env %s", k) }
	return v
}
