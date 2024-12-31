docker stop drunkgames
docker rm drunkgames
docker build -t drunkgames .
docker run -d -p 8080:8080 --name drunkgames drunkgames