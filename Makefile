.PHONY: 

# install buf & grpcurl
install:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

buf:
	buf lint
	buf dep update
	buf generate

#run all services
# run-services:
# 	${MAKE} -C ./services/svc-order

# TODO
# test all service unit test
# copy all service .env.example
# docker compose up
# docker compose down

# run-services:
# 	$(MAKE) -C ./services/svc-inventory run
# 	$(MAKE) -C ./services/svc-order run

