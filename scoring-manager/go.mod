module github.com/KlassnayaAfrodita/github-user-score/scoring-manager

go 1.23

toolchain go1.23.6

require (
	//github.com/KlassnayaAfrodita/github-user-score/collector v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.6.0
	github.com/jackc/pgx/v5 v5.7.4
	github.com/robfig/cron/v3 v3.0.1
	github.com/segmentio/kafka-go v0.4.47
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
)

require github.com/KlassnayaAfrodita/github-user-score/collector v0.0.0-20250514093611-6f5fe103f866

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/KlassnayaAfrodita/github-user-score/collector => ../collector
