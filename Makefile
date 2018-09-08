build: hub app_rise stompy gob app_sink

app_sink:	
	go build ./cmd/app_sink

hub:
	go build ./cmd/hub

app_rise:
	go build ./cmd/app_rise

stompy:
	go build ./cmd/stompy

gob:
	go build ./cmd/gob

clean:
	rm -f app_sink
	rm -f hub
	rm -f app_rise
	rm -f stompy
	rm -f gob

rebuild: clean build
