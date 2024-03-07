package internal

import "time"

type Config struct {
	Address    string        `env:"ADDRESS" envDefault:"localhost:9000"`
	Expiration time.Duration `env:"EXPIRATION" envDefault:"30m"`
	GCInterval time.Duration `env:"GC_INTERVAL" envDefault:"1m"`
	ChunkCount int           `env:"CHUNK_COUNT" envDefault:"2048"`
}
