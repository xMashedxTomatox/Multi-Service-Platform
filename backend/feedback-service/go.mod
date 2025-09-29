module github.com/xmashedxtomatox/feedback-service

go 1.25.1

require (
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.11.1
	github.com/xmashedxtomatox/shared v0.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/xmashedxtomatox/shared => ../shared
