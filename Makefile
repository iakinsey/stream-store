default: clean
	go build .

suru: clean
	go build .
	echo '{"executables": [{"uuid": "af9368ed-7552-44f1-9c71-54b2339ce24f", "name": "stream-store", "count": 1, "executable": "./stream-store"}]}' > config.json
	zip stream-store-suru.zip config.json stream-store
	rm -f stream-store
	rm -f config.json

docker: clean
	docker build . -t stream-store
	docker run -d stream-store
	docker cp `docker ps | grep stream-store | tr -s ' ' | cut -d " " -f 1`:/build/stream-store .
	docker kill `docker ps | grep stream-store | tr -s ' ' | cut -d " " -f 1`

clean:
	sh -c "docker kill `docker ps | grep stream-store | tr -s ' ' | cut -d " " -f 1`" || :
	rm -f stream-store
	rm -f stream-store-suru.zip
	rm -f config.json
