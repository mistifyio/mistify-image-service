PREFIX := /usr
SBIN_DIR=$(PREFIX)/sbin
CONF_DIR := /etc/mistify-image-service

cmd/mistify-image-service/mistify-image-service: cmd/mistify-image-service/main.go
	cd $(dir $<) && \
	go get && \
	go build

$(SBIN_DIR)/mistify-image-service: cmd/mistify-image-service/mistify-image-service
	install -D $< $(DESTDIR)$@

$(CONF_DIR)/config.json: cmd/mistify-image-service/config.json
	install -D -m 0444 $< $(DESTDIR)$@

clean:
	cd cmd/mistify-image-service && \
	go clean

install: \
	$(CONF_DIR)/config.json \
	$(SBIN_DIR)/mistify-image-service
