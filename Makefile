KEY:=~/secret/$(PROJECT)-pgadaptor.json
INSTANCE:=test-instance
DATABASE_NAME:=game
CONTAINER_NAME:=gorm-pgadaptor

.PHONY: start-pgadaptor
start-pgadaptor:
	docker run -d --rm -p 15432:5432 --name $(CONTAINER_NAME) \
    -v $(KEY):/acct_credentials.json \
    gcr.io/cloud-spanner-pg-adapter/pgadapter:latest \
    -p $(PROJECT) -i $(INSTANCE) -d $(DATABASE_NAME) \
    -c /acct_credentials.json -q -x

.PHONY: stop-pgadapter
stop-pgadaptor:
	docker stop $(CONTAINER_NAME)

.PHONY: make-schemas
make-schemas:
	psql -h localhost -p 15432 musics < ./schemas/ddl.sql 
