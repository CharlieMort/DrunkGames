docker rm -f drunkgames-go-web-app-1
docker rm -f drunkgames-letsencrypt-nginx-proxy-companion-1
docker rm -f drunkgames-nginx-proxy-1
docker compose -f nginx-proxy-compose.yaml up -d
docker compose -f go-app-compose.yaml up -d