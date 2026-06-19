
build:
	@echo "Building for your architecture"
	go build -o bin/our .

install: build
	@echo Installing builded package
	cp bin/ssl-manager /usr/bin/our
