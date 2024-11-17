# protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    routeguide/route_guide.proto

gen:
	protoc --go-grpc_out=. --go_out=. ./proto/*.proto

server: gen server/main.go
	go run server/main.go

open:
	# Open the grpcui in browser
	grpcui -cert certs/certfile.crt -key certs/keyfile.key -cacert certs/rootca.crt localhost:9000

client: gen client/main.go
	go run client/main.go

.PHONY: gen server open client