default: clean
	go build .

suru: clean
	go build .
	echo '{"executables": [{"uuid": "af9368ed-7552-44f1-9c71-54b2339ce24f", "name": "stream-store", "count": 1, "executable": "./stream-store"}]}' > config.json
	zip stream-store-suru.zip config.json stream-store
	rm -f stream-store
	rm -f config.json

clean:
	rm -f stream-store
	rm -f stream-store-suru.zip
	rm -f config.json