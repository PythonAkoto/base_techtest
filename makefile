.PHONY: start stop restart logs

# Start the containers in detached mode
start:
	docker-compose up -d

# Stop the containers
stop:
	docker-compose down

# Restart the containers
restart:
	docker-compose down && docker-compose up -d

# View logs
logs:
	docker-compose logs -f
