dev:
	fresh

prod:
	go build -o bin/blogchi.exe ./blogchi
	bin/blogchi.exe -env=production -port=80

prod_unix:
	go build -o bin/blogchi ./blogchi
	bin/blogchi -env=production -port=80
