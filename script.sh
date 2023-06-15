
sudo docker run --network=host todolist
docker rm -f $(docker ps -aq)
docker network inspect my-network
docker network create my-network
docker build --no-cache -t todolist .
docker create --name my-mysql -e MYSQL_ROOT_PASSWORD=Pastibisa -e MYSQL_DATABASE=Gin_todo -p 3307:3306 --network my-network mysql:8.0.33
docker run -d --name my-mysql -e MYSQL_ROOT_PASSWORD=Pastibisa -e MYSQL_DATABASE=Gin_todo -p 3307:3306 --network my-network mysql:8.0.33
docker create --name my-app -p 8080:8080 --network my-network my-app-image
docker run -d --name my-app -p 8080:8080 --network my-network my-app-image