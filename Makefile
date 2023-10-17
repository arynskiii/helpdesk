mysql:
	sudo docker run -d -e MYSQL_ROOT_PASSWORD=qwerty -e MYSQL_DATABASE=remote -e MYSQL_USER=mysql1 -e MYSQL_PASSWORD=qwerty -p 3306:3306 --name=mysql_container mysql
mysqlterminal:
	sudo docker exec -it YOUR_CONTAINER_ID|YOUR_CONTAINER_NAME mysql -u mysql -p