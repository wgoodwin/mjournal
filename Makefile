# If PREFIX isn't provided, we check for /usr/local and use that if it exists.
# Otherwise we fall back to using /usr

LOCAL != test -d $(DESTDIR)/usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d $(DESTDIR)/usr/local && echo "/local" || echo ""
PREFIX ?= /usr$(LOCAL)

build:
	#  Build directives here

install:
	# Install directives here

uninstall:
	# Uninstall directives here


