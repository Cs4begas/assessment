start-project:
	DATABASE_URL=postgres://vmgwwcfs:h1SnYCJMWCfmiJ_vbx3tE-ZN3H0rcMTk@rosie.db.elephantsql.com/vmgwwcfs PORT=2565 go run server.go

unit-test:
	go test -v -tags unit ./...

integration-test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests

coverage-test:
	go test -tags unit -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out