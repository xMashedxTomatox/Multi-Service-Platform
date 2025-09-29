module github.com/xmashedxtomatox/auth-service

go 1.25.1

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/gorilla/mux v1.8.1
	github.com/redis/go-redis/v9 v9.15.0
	golang.org/x/crypto v0.37.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/lib/pq v1.10.9 // indirect
)

require github.com/xmashedxtomatox/shared v0.0.0

replace github.com/xmashedxtomatox/shared => ../shared
