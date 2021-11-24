.PHONY: all client server

all: client server

aarsdir:
	ls aars || mkdir aars

client: aarsdir client/*
	gomobile bind -v -o aars/client.aar ftp/client ftp/fm

server: aarsdir server/*
	gomobile bind -v -o aars/server.aar ftp/server

clean:
	rm -r aars || true
