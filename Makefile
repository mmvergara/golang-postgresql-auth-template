db:
	docker compose up -d;
	psql -h localhost -p 5432 -U root -d gpa;


