go:
	sudo docker compose up -d

conn:
	sudo docker compose -f docker-compose.yml exec db mongosh -u root -p example

rm:
	sudo docker rm -vf $$(sudo docker ps -aq)
	sudo docker rmi -f $$(sudo docker images -aq)

rm_f: rm
	sudo docker volume rm $$(sudo docker volume ls -f dangling=true -q)
