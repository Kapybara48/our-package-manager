
PREFIX := /usr/local
BINDIR := $(PREFIX)/bin
LIBDIR := $(PREFIX)/lib
SHRDIR := $(PREFIX)/share
DOCDIR := $(PREFIX)/share/doc
MANDIR := $(PREFIX)/share/man


build:
	@echo "Building for your architecture"
	go build -o bin/our .

install: build
	@echo Installing builded package
	install -Dm755 bin/our -t '$(DESTDIR)$(BINDIR)'

uninstall:
	@echo Removing binary files
	rm -rf '$(DESTDIR)$(BINDIR)/our'

