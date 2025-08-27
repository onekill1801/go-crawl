# migrate -path ./migrations -database "mysql://root:your_root_password@tcp(192.168.1.6:5306)/test" drop -f
migrate -path ./migrations -database "mysql://root:your_root_password@tcp(192.168.1.6:5306)/test" up
# migrate -path ./migrations -database "mysql://root:your_root_password@tcp(192.168.1.6:5306)/test" down

# go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest add variable to PATH
