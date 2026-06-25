
build:
	@echo "Building for your architecture"
	go build -o bin/our .

install: build
	@echo Installing builded package
	rm -rf /usr/sbin/our
	cp ./bin/our /usr/sbin/our

uninstall:
	@echo Removing binary files
	rm -rf /usr/sbin/our

