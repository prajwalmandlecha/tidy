BINARY=tidy
INSTALL_PATH=/usr/local/bin/$(BINARY)
CONFIG_PATH=$(HOME)/.config/tidy/config.yaml

build:
	go build -o $(BINARY) .

install: build
	sudo mv $(BINARY) $(INSTALL_PATH)
	systemctl --user daemon-reload
	systemctl --user restart $(BINARY)
	@echo "Run 'tidy init' to generate your config file"

uninstall:
	systemctl --user stop $(BINARY)
	systemctl --user disable $(BINARY)
	sudo rm -f $(INSTALL_PATH)

logs:
	journalctl --user -u $(BINARY) -f

status:
	systemctl --user status $(BINARY)

test:
	go test -v ./...

.PHONY: build install uninstall logs status test