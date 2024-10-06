build:
	go build -o build/pills-cron ./cmd/cron && \
	go build -o build/pills-bot ./

start:
	make build && \
	pm2 start ecosystem_config.yml

rebuild:
	rm -R build && \
	make build && \
	make restart

restart:
	pm2 delete pills-cron pills-bot && \
	pm2 start ecosystem_config.yml