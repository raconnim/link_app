run_in_memory:
	docker-compose build app_in_memory
	docker-compose up app_in_memory

run_in_db:
	docker-compose build app
	docker-compose up app