.PHONY: test-e2e test coverage

test-e2e:
	go test -v -count=1 -parallel=4 -coverpkg=av-merch-shop/internal/... -coverprofile=cov.out ./tests/e2e/...

test:
	go test -v -count=1 -parallel=4 -coverpkg=av-merch-shop/internal/... -coverprofile=cov.out ./... ./tests/e2e/...
	go tool cover -func=cov.out

coverage:
	go tool cover -html=cov.out
