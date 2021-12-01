# deployments
test-server:
	docker-compose -f ./deployments/docker-compose-development.yml build --no-cache --force-rm
	docker-compose -f ./deployments/docker-compose-development.yml up

test-server-down:
	docker-compose -f ./deployments/docker-compose-development.yml down

prod-server:
	docker-compose -f ./deployments/docker-compose-production.yml build --no-cache --force-rm
	docker-compose -f ./deployments/docker-compose-production.yml up
	
prod-server-down:
	docker-compose -f ./deployments/docker-compose-production.yml down

# testing cmds
run-integration-tests:
	go test -v -cover -tags integration ./...

run-unit-tests:
	go test -v -cover -tags unit ./...

run-all-tests:
	go test -v -cover -tags=unit,integration ./...
	# run-unit-tests run-integration-tests


# mocking section
update-database-mocks:
	mockery --name=MongoQueries --structname=MongoQueriesMock --filename=mgdb_queries.go --recursive --output ./internal/mocks

update-service-mocks:
	mockery --name=ApiServiceInterface --structname=ApiServicesMock --filename=services_api.go --recursive --output ./internal/mocks
	mockery --name=UserAuthServicesInterface --structname=UserAuthServicesInterfaceMock --filename=service_auth.go --recursive --output ./internal/mocks

# remock ALL tracked mock with updates
update-mocks: update-service-mocks update-database-mocks
	