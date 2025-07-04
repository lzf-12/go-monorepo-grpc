.PHONY: 

# install buf & grpcurl
install:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

buf:
	buf lint
	buf dep update
	buf generate

# TODO
# test all service unit test
# copy all service .env.example
# docker compose up
# docker compose down
