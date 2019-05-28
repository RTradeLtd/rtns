# Rebuild generate code
COUNTERFEITER=go run github.com/maxbrunsfeld/counterfeiter/v6
.PHONY: gen
gen:
	@echo "===================    regenerating code    ==================="
	$(COUNTERFEITER) -o ./mocks/kaas.mock.go \
		github.com/RTradeLtd/grpc/krab.ServiceClient
	$(COUNTERFEITER) -o ./mocks/namesys.mock.go \
		github.com/ipfs/go-ipfs/namesys.NameSystem