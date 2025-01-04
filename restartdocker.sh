docker compose -f nginx-proxy-compose down
docker compose -f go-app-compose.yaml down
docker compose -f nginx-proxy-compose.yaml build --no-cache
docker compose -f go-app-compose.yaml build --no-cache
docker compose -f nginx-proxy-compose.yaml up -d
docker compose -f go-app-compose.yaml up -d