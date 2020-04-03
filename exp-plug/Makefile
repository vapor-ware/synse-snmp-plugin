build:
	GOOS=linux GOARCH=amd64 go build -o expplug .
	docker build -t etd/expplugin .

run:
	#docker run --name expplug -p 5001:5001 etd/expplugin
	docker-compose up -d