
build:
	@echo "Building for your architecture"
	go build -o bin/our .

install: build
	@echo Installing builded package
	mv /usr/bin/our /usr/bin/our.old
	cp ./bin/our /usr/bin/our
