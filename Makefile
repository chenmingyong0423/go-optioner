# 单元测试
.PHONY: ut
ut:
	@go test -race ./...

.PHONY: setup
setup:
	@sh ./script/setup.sh

.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	@sh ./script/fmt.sh

.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: check
check:
	@$(MAKE) --no-print-directory fmt
	@$(MAKE) --no-print-directory tidy

.PHONY: e2e_up
e2e_up:
	docker-compose -f script/integration_test_compose.yml up -d

.PHONY: e2e_down
e2e_down:
	docker-compose -f script/integration_test_compose.yml down -v