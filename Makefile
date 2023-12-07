default:
	go build -o ./dist/netcat-test main.go
run:
	rm output.csv||true
	./dist/netcat-test -f ./test.csv -o output.csv
	cat output.csv

build-ubuntu-22.04:
	docker build -t netcat-tester:ubuntu-22.04 -f Dockerfile.ubuntu-22.04 . 
	docker container rm -f builder||true
	docker run --name=builder -d netcat-tester:ubuntu-22.04 sleep infinity 
	docker cp builder:/netcat-tester netcat-tester-ubuntu-22
	mv netcat-tester-ubuntu-22 ./dist/

build-fedora-32:
	docker build -t netcat-tester:fedora-32 -f Dockerfile.fedora-32 . 
	docker container rm -f builder||true
	docker run --name=builder -d netcat-tester:fedora-32 sleep infinity 
	docker cp builder:/netcat-tester netcat-tester-fedora-32
	mv netcat-tester-fedora-32 ./dist/

build-all-docker:
	make build-ubuntu-22.04
	make build-fedora-32