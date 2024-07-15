# vpn
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr/local
MANPREFIX ?= $(PREFIX)/share/man
GO ?= go
GOFLAGS ?=
RM ?= rm -f

all: vpn

vpn:
	$(GO) build $(GOFLAGS)
	scdoc < vpn.1.scd > vpn.1

clean:
	$(RM) vpn
	$(RM) vpn.1

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f vpn $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/vpn
	mkdir -p $(DESTDIR)$(MANPREFIX)/man1
	cp -f vpn.1 $(DESTDIR)$(MANPREFIX)/man1/vpn.1
	chmod 644 $(DESTDIR)$(MANPREFIX)/man1/vpn.1

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/vpn
	$(RM) $(DESTDIR)$(MANPREFIX)/man1/vpn.1

.DEFAULT_GOAL := all

.PHONY: all vpn clean install uninstall
