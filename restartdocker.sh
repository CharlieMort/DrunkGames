docker compose -f go-app-compose.yaml down
docker compose -f go-app-compose.yaml build --no-cache
docker compose -f go-app-compose.yaml up -d