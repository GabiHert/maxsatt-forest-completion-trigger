module github.com/GabiHert/maxsatt-forest-completion-trigger

go 1.25.1

require (
	github.com/alicebob/miniredis/v2 v2.35.0
	github.com/aws/aws-lambda-go v1.50.0
	github.com/aws/aws-sdk-go-v2 v1.39.6
	github.com/aws/aws-sdk-go-v2/config v1.31.20
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.20
	github.com/aws/aws-sdk-go-v2/service/batch v1.58.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.52.3
	github.com/aws/aws-sdk-go-v2/service/rekognition v1.51.8
	github.com/aws/aws-sdk-go-v2/service/s3 v1.90.2
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.39.10
	github.com/aws/aws-sdk-go-v2/service/ses v1.34.8
	github.com/aws/aws-sdk-go-v2/service/sns v1.39.2
	github.com/aws/aws-sdk-go-v2/service/sqs v1.42.12
	github.com/cucumber/godog v0.15.1
	github.com/go-playground/validator/v10 v10.28.0
	github.com/gofiber/fiber/v2 v2.52.9
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/lsgndln/dd-trace-go v1.3.0
	github.com/parquet-go/parquet-go v0.25.1
	github.com/redis/go-redis/v9 v9.16.0
	golang.org/x/crypto v0.45.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.0
)

require (
	github.com/DataDog/appsec-internal-go v1.0.0 // indirect
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.45.0-rc.1 // indirect
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.45.0-rc.1 // indirect
	github.com/DataDog/datadog-go/v5 v5.1.1 // indirect
	github.com/DataDog/go-libddwaf v1.2.0 // indirect
	github.com/DataDog/go-tuf v0.3.0--fix-localmeta-fork // indirect
	github.com/DataDog/sketches-go v1.2.1 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.24 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.13 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.32.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.40.2 // indirect
	github.com/aws/smithy-go v1.23.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denisenkom/go-mssqldb v0.12.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.10 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20221118152302-e6195bd50e26 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/outcaste-io/ristretto v0.2.1 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.5.0 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go4.org/intern v0.0.0-20211027215823-ae77deb06f29 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20220617031537-928513b29760 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231120223509-83a465c0220f // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	inet.af/netaddr v0.0.0-20220811202034-502d2d690317 // indirect
)
