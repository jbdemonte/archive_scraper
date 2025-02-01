APP := archive_scraper

build:
	go build -o ${APP}

build_x86_64:
	GOOS=linux GOARCH=amd64 go build -o ${APP}

build_armv7l:
	GOOS=linux GOARCH=arm GOARM=7 go build -o ${APP}

build_aarch64:
	GOOS=linux GOARCH=arm64 go build -o ${APP}

build_mipsel:
	GOOS=linux GOARCH=mipsle go build -o ${APP}