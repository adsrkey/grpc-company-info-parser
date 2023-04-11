gen:
	rm -f parser/parserpb/*.go
	protoc --proto_path=parser \
	--go_out=parser/parserpb --go_opt paths=source_relative \
   	--go-grpc_out=parser/parserpb  --go-grpc_opt paths=source_relative \
	--grpc-gateway_out=parser/parserpb --grpc-gateway_opt paths=source_relative \
	./parser/*.proto

doc:
	protoc -I parser --swagger_out=logtostderr=true:. parser/parser.proto

up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

up_build:
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

down:
	@echo "Stopping docker compose..."
	docker-compose down -d
	@echo "docker-compose down: Done!"