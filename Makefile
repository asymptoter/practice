# 定義變量
ENV ?= local
image_tag = $(shell date +"%Y%m%d%H%M%S")

# 根據 ENV 變量選擇要執行的目標
ifeq ($(ENV), local)
TARGET = local
else ifeq ($(ENV), prod)
TARGET = prod
else
$(error Unknown environment: $(ENV))
endif

.PHONY: all local prod

# 默認目標，根據選擇的環境執行相應的目標
all: $(TARGET)

# local 環境下執行 docker compose up
local:
	docker compose build
	docker compose up && goose postgres "user=asymptoter password=password dbname=practice sslmode=disable" up

# prod 
prod:
	docker tag server:latest asia-east1-docker.pkg.dev/practice-backend-423606/practice-backend/server:latest
	docker push asia-east1-docker.pkg.dev/practice-backend-423606/practice-backend/server:latest
