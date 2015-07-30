build: *.go
		go build -o build/helios-github

gin:
		@gin -a 8888 -p 8989
