.PHONY: build build-web build-go dev-api dev-web docker clean

# Full production build: Angular → Go binary with embedded SPA
build: build-web build-go

build-web:
	cd web && npm install && npx ng build --configuration=production

build-go: build-web
	go build -o gesitr .

# Development: run these in separate terminals
dev-api:
	DEV=true go run .

dev-web:
	cd web && npx ng serve

docker:
	docker build -t gesitr .

clean:
	rm -f gesitr gesitr.db
	rm -rf web/dist
