.PHONY: all client server fm

all: client server fm

aarsdir:
	ls aars || mkdir aars

client: aarsdir client/*
	gomobile bind -v -o aars/client.aar ftp/client

server: aarsdir server/*
	gomobile bind -v -o aars/server.aar ftp/server

fm: aarsdir fm/*
	gomobile bind -v -o aars/fm.aar ftp/fm

clean:
	rm -r aars || true
