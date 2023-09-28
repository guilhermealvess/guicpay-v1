deps:
	go mod download
	go mod tidy

sql-generate:
	docker run --rm -v /home/guilherme.alves/projects/chalanges/guicpay:/src -w /src sqlc/sqlc generate