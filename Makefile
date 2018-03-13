dev:
	fresh

prod:
	go build -o bin/blogchi.exe ./blogchi
	bin/blogchi.exe -env=production -port=80
