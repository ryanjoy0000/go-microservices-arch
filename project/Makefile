FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp

# ----------------------------
# up : start all containers in background without forcing build
up:
	@echo "Starting all services (except frontend) in background..."
	docker-compose up -d

# up_build: stops docker-compose, builds all projects and starts docker-compose
up_build: build_broker build_auth
	@echo "Stopping docker containers if any..."
	docker-compose down
	@echo "Building & starting all services (except frontend) in background..."
	docker-compose up --build

# down : stop docker-compose
down:
	@echo "Stopping docker-compose..."
	docker-compose down


# ------------------------------

# build_broker: Build the broker service -> binary as a linux executable
build_broker:
	@echo "Building the broker service -> binary..."
	cd ../broker-service && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${BROKER_BINARY} ./cmd/api

# build_auth: Build the auth service -> binary as a linux executable
build_auth:
	@echo "Building the auth service -> binary..."
	cd ../auth-service && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${AUTH_BINARY} ./cmd/api

# build_front: Build the front end
build_front:
	@echo "Building the front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web

# start_front: Build & Start the front end
start_front: build_front
	@echo "Starting front end..."
	cd ../front-end && ./${FRONT_END_BINARY}

# stop_front: Stop the front end
stop_front: 
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"