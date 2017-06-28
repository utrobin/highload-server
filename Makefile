all:
	go build -i highload/src/main
	mv main httpd