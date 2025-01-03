docker stop drunkgames-go-web-app-1
docker stop drunkgames-letsencrypt-nginx-proxy-companion-1
docker stop drunkgames-nginx-proxy-1
docker compose -f nginx-proxy-compose.yaml up -d
docker compose -f go-app-compose.yaml up -d