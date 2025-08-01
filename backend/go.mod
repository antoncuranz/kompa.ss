module travel-planner

go 1.24

tool (
	github.com/pressly/goose/v3/cmd/goose
	github.com/sqlc-dev/sqlc/cmd/sqlc
	github.com/swaggo/swag/cmd/swag
	go.uber.org/mock/mockgen
)

require (
	cloud.google.com/go v0.121.4
	github.com/ansrivas/fiberprometheus/v2 v2.11.0
	github.com/caarlos0/env/v11 v11.3.1
	github.com/go-playground/validator/v10 v10.26.0
	github.com/goccy/go-json v0.10.5
	github.com/gofiber/fiber/v2 v2.52.8
	github.com/gofiber/swagger v1.1.1
	github.com/jackc/pgx/v5 v5.7.5
	github.com/pressly/goose/v3 v3.24.3
	github.com/rs/zerolog v1.34.0
	github.com/swaggo/swag v1.16.4
)

require (
	cel.dev/expr v0.24.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/ClickHouse/ch-go v0.66.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.35.0 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coder/websocket v1.8.13 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/cubicdaiya/gonp v1.0.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/elastic/go-sysinfo v1.15.3 // indirect
	github.com/elastic/go-windows v1.0.2 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/mfridman/xflag v0.1.0 // indirect
	github.com/microsoft/go-mssqldb v1.8.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pganalyze/pg_query_go/v6 v6.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pingcap/errors v0.11.5-0.20240311024730-e056997136bb // indirect
	github.com/pingcap/failpoint v0.0.0-20240528011301-b51a646c7c86 // indirect
	github.com/pingcap/log v1.1.0 // indirect
	github.com/pingcap/tidb/pkg/parser v0.0.0-20250531022214-e7b038b99132 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.64.0 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/riza-io/grpc-go v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/sqlc-dev/sqlc v1.29.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/swaggo/files/v2 v2.0.2 // indirect
	github.com/tetratelabs/wazero v1.9.0 // indirect
	github.com/tursodatabase/libsql-client-go v0.0.0-20240902231107-85af5b9d094d // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.62.0 // indirect
	github.com/vertica/vertica-sql-go v1.3.3 // indirect
	github.com/wasilibs/go-pgquery v0.0.0-20250409022910-10ac41983c07 // indirect
	github.com/wasilibs/wazero-helpers v0.0.0-20250123031827-cd30c44769bb // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	github.com/ydb-platform/ydb-go-genproto v0.0.0-20250519101544-1f330d77b70f // indirect
	github.com/ydb-platform/ydb-go-sdk/v3 v3.108.5 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/exp v0.0.0-20250531010427-b6e5de432a8b // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/grpc v1.73.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.1 // indirect
	modernc.org/libc v1.65.8 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.37.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
