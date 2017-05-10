default:
	go install
race:
	go install --race
build:
	cd $(GOPATH)/bin; \
		./Cross\ Compile\ C.sh gofilesync gofilesync
