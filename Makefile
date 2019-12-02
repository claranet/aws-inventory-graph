#!make

up:
	docker-compose -f docker-compose.yml up -d

stop: 
	docker-compose -f docker-compose.yml stop

rm: 
	docker-compose -f docker-compose.yml rm