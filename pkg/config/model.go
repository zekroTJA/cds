package config

var Default = Config{
	Address: "0.0.0.0:80",
	Logging: Logging{
		Level: "info",
	},
}

type Config struct {
	Address string
	Logging Logging
	Stores  []Store
}

type Logging struct {
	Level string
}

type StoreType string

const (
	StoreTypeLocal = StoreType("local")
	StoreTypeS3    = StoreType("s3")
)

type Store struct {
	Type         StoreType
	Entrypoint   string
	Listable     bool
	CacheControl string

	// Shared Params
	Path string

	// S3 Params
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	Secure    bool
	Bucket    string
}
