build:
	docker build --pull --rm -f "Dockerfile" -t ruiblaese/projeto-golang:latest "."

run:
	docker run -p "8080:8080" ruiblaese/projeto-golang 

push:	
	docker build --pull --rm -f "Dockerfile" -t ruiblaese/projeto-golang:latest "."
	docker tag ruiblaese/projeto-golang:latest ruiblaese/projeto-golang:0.3.0	
	docker push ruiblaese/projeto-golang:0.3.0	

redis: 
	docker run --name redis -p 6379:6379 -d redis