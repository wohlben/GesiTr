.PHONY: build build-web build-go dev dev-api dev-web docker clean dolt generate seed test test-go test-web test-e2e lint lint-go lint-web update-screenshots update-screenshots-web update-screenshots-e2e

# Generate TypeScript types from Go structs
generate:
	go run github.com/gzuidhof/tygo@latest generate

# Full production build: generate types → Angular → Go binary with embedded SPA
build: generate build-web build-go

build-web:
	cd web && npm install && npx ng build --configuration=production

build-go: build-web
	go build -o gesitr .

# Development: start both API and web servers
dev:
	$(MAKE) -j2 dev-api dev-web

dev-api:
	DEV=true go run .

dev-web: generate
	cd web && npx ng serve

dolt:
	dolt sql-server --host 127.0.0.1 --port 3307 --data-dir .beads/dolt

# Lint and format check
lint: lint-go lint-web

lint-go:
	@test -z "$$(gofmt -l .)" || (echo "Go files not formatted:" && gofmt -l . && exit 1)

lint-web:
	cd web && npm run lint && npm run format:check

# Run all tests
test: lint test-go test-web test-e2e

test-go:
	go test ./...

test-web:
	cd web && npm test

test-e2e:
	cd web && npx ng e2e

# Update screenshot baselines
update-screenshots: update-screenshots-web update-screenshots-e2e

update-screenshots-web:
	find web/src/app -path '*__screenshots__*' -name '*.png' -delete
	cd web && npx ng run web:test-screenshot || true
	cd web && npx ng run web:test-screenshot

update-screenshots-e2e:
	cd web && npx ng e2e --update-snapshots

docker:
	docker build -t gesitr .

seed:
	rm -f gesitr.db
	go run ./cmd/seed

clean:
	rm -f gesitr gesitr.db
	rm -rf web/dist
	rm -f web/src/app/generated/models.ts
