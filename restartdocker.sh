docker compose down drunkgames-go-web-app-1
docker compose down drunkgames-letsencrypt-nginx-proxy-companion-1
docker compose down drunkgames-nginx-proxy-1
docker compose -f nginx-proxy-compose.yaml build --no-cache
docker compose -f go-app-compose.yaml build --no-cache
docker compose -f nginx-proxy-compose.yaml up -d
docker compose -f go-app-compose.yaml up -d