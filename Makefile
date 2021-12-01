build: frontend.build backend.build

push: frontend.push backend.push

frontend.push: frontend.build
	docker push ghcr.io/tobiaskohlbau/unturned-admin/frontend:latest

backend.push: backend.build
	docker push ghcr.io/tobiaskohlbau/unturned-admin/backend:latest

frontend.build:
	docker build -t ghcr.io/tobiaskohlbau/unturned-admin/frontend:latest ./web

backend.build:
	docker build -t ghcr.io/tobiaskohlbau/unturned-admin/backend:latest .
