tidy:
	go mod tidy

run: 
	go run main.go

build:
	docker build -t mindspace-backend .

# docker-compose-dev
up: 
	docker-compose up -d

down:
	docker-compose down

# docker-compose-prod
up-prod:
	docker-compose -f docker-compose-prod.yaml up -d

down-prod:
	docker-compose -f docker-compose-prod.yaml down

exec-pg:
	docker exec -it mindspace-db psql -U postgres -d mindspace