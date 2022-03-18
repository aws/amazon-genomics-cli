PACKAGES := client 

.PHONY: format test build init $(PACKAGES)

all: format test build

format:
	target=format $(MAKE) $(PACKAGES)

test:
	target=test $(MAKE) $(PACKAGES)

build: 
	target=build $(MAKE) $(PACKAGES)

init:
	target=init $(MAKE) $(PACKAGES)
	
$(PACKAGES):
	(cd $@; $(MAKE) $(target))

