default:
	go build -o ./dist/netcat-test main.go
run:
	rm output.csv||true
	./dist/netcat-test -f ./test.csv -o output.csv
	cat output.csv
