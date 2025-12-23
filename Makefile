TEST_THRESHOLD_COVERAGE := 85.0
MUTATION_THRESHOLD_COVERAGE := 89.0
MUTATION_THRESHOLD_EFFICACY=100

.PHONY: run build-binary validate

run:
	go run ./cmd/local/main.go

build-binary:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o ${app_name} -ldflags="-s -w" ./cmd/app/main.go

build-binary-local:
	CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap -ldflags="-s -w" ./cmd/app/main.go

validate:
	@echo "Sort imports..."
	make sort-imports
	@echo "Fix struct alignment..."
	make fix-struct-alignment
	@echo "Running dependency check..."
	make dependency-check
	@echo "Running unit tests..."
	make tests
	@echo "Validating test coverage..."
	make validate-coverage
	@echo "Running mutation tests..."
	make mutation-test
	@echo "Validating mutation coverage..."
	make validate-mutation-coverage
	@echo "Validating test efficacy..."
	make validate-test-efficacy

tests:
	@echo "Running tests..."
	go test ./test/... -coverpkg=./internal/... -coverprofile=coverage.out ./... -p=100
	go tool cover -func coverage.out
	@echo "Tests completed!"

tests-with-tag:
	@echo "Running tests..."
	go test ./test/... -v -coverprofile=coverage.out -coverpkg=./internal/... -p 100 --scenarios=@all
	go tool cover -html=coverage.out -o coverage.html
	@echo "Tests completed!"

validate-coverage:
	@echo "Validating coverage..."
	@coverage=$$(go tool cover -func=coverage.out | grep '^total:' | awk '{print $$NF}' | sed 's/%//'); \
	threshold=$(TEST_THRESHOLD_COVERAGE); \
	if (( $$(echo "$$coverage < $$threshold" | bc -l) )); then \
		echo "❌ Test coverage ($$coverage%) is below the required threshold ($$threshold%)."; \
		exit 1; \
	else \
		echo "✅ Test coverage meets the required threshold ($$coverage% >= $$threshold%)."; \
	fi
	@echo "Coverage validation completed!"

mutation-test:
	@echo "Installing Gremlins..."
	go get github.com/go-gremlins/gremlins/cmd/gremlins
	go install github.com/go-gremlins/gremlins/cmd/gremlins
	go mod tidy
	@echo "Gremlins installed successfully!"
	@echo "Running mutation tests..."
	gremlins unleash \
              --silent=false \
              --integration=true \
              --dry-run=false \
              --tags="" \
              --output="" \
              --workers=4 \
              --test-cpu=2 \
              --timeout-coefficient=0 \
              --threshold-efficacy=90 \
              --threshold-mcover=50 \
              --exclude-files "^cmd/" \
              --exclude-files "^pkg/" \
              --exclude-files "^test/" \
              --exclude-files "^config/" \
              --exclude-files "^build/" \
              --coverpkg "./internal/..." \
              --arithmetic-base=true \
              --conditionals-boundary=true \
              --conditionals-negation=true \
              --increment-decrement=true \
              --invert-assignments=true \
              --invert-bitwise=true \
              --invert-bwassign=true \
              --invert-negatives=true \
              --invert-logical=true \
              --invert-loopctrl=true \
              --remove-self-assignments=true | tee mutation_results.log
	@echo "Mutation tests completed!"

validate-mutation-coverage:
	@echo "Validating mutation coverage..."
	coverage=$$(grep -E 'Mutator coverage: [0-9]+(\.[0-9]+)?%' mutation_results.log | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	threshold=$(MUTATION_THRESHOLD_COVERAGE); \
	result=$$(echo "$$coverage >= $$threshold" | bc); \
	if [ "$$result" -eq 0 ]; then \
		echo "❌ Mutation test coverage ($$coverage%) is below the required threshold ($$threshold%)."; \
		exit 1; \
	else \
		echo "✅ Mutation test coverage meets the required threshold ($$coverage%)/($$threshold%)."; \
	fi
	@echo "Mutation coverage validation completed!"

validate-test-efficacy:
	@echo "Validating test efficacy..."
	efficacy=$$(grep -E 'Test efficacy: [0-9]+(\.[0-9]+)?%' mutation_results.log | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	threshold=$(MUTATION_THRESHOLD_EFFICACY); \
	result=$$(echo "$$efficacy >= $$threshold" | bc); \
	if [ "$$result" -eq 0 ]; then \
		echo "❌ Mutation test coverage ($$efficacy%) is below the required threshold ($$threshold%)."; \
		exit 1; \
	else \
		echo "✅ Mutation test efficacy meets the required threshold ($$efficacy%)/($$threshold%)."; \
	fi
	@echo "Mutation coverage validation completed!"

dependency-check:
	@echo "Installing Nancy..."
	go install github.com/sonatype-nexus-community/nancy@latest
	go mod tidy
	@echo "Nancy installed successfully!"
	@echo "Running dependency check..."
	@if [ -z "$$NANCY_USERNAME" ] || [ -z "$$NANCY_TOKEN" ]; then \
		echo "Error: NANCY_USERNAME and NANCY_TOKEN environment variables must be set"; \
		exit 1; \
	fi
	go list -json -deps ./... | nancy sleuth --username "$$NANCY_USERNAME" --token "$$NANCY_TOKEN" | tee nancy-report.log
	@echo "Dependency check completed!"

autofix-dependency-check:
	@echo "Running dependency check fixes..."
	findings_file="nancy-report.log"; \
	grep -Eo 'pkg:golang/[a-zA-Z0-9._/-]+@[v0-9]+\.[0-9]+\.[0-9]+' $$findings_file | while read -r line; do \
		module=$$(echo $$line | sed 's/pkg:golang\///' | awk -F'@' '{print $$1}'); \
		echo "Updating $$module to the latest version..."; \
		go get -u "$$module"; \
	done; \
	go mod tidy
	@echo "Dependency check fix completed!"

sort-imports:
	@echo "Sorting imports..."
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .
	@echo "Imports sorted!"

fix-struct-alignment:
	@echo "Installing betteralign..."
	@go install github.com/dkorunic/betteralign/cmd/betteralign@latest
	@echo "Fixing struct alignment issues..."
	@GOBIN=$$(go env GOBIN); \
	if [ -z "$$GOBIN" ]; then GOBIN="$$(go env GOPATH)/bin"; fi; \
	while true; do \
		"$$GOBIN/betteralign" -apply ./...; \
		exit_code=$$?; \
		if [ $$exit_code -eq 0 ]; then \
			break; \
		elif [ $$exit_code -eq 3 ]; then \
			echo "Running betteralign again to apply remaining fixes..."; \
		else \
			echo "Error: betteralign failed with exit code $$exit_code"; \
			exit $$exit_code; \
		fi; \
	done
	@echo "✅ Struct alignment fixed!"
