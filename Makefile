
prybar: main.go plugins/python2.so plugins/python3.so plugins/ruby.so plugins/lua.so
	go build -o prybar

plugins/%.so: languages/%/*.go
	cd languages/$* && CGO_LDFLAGS_ALLOW=".*" go build -buildmode=plugin -o ../../$@

clean:
	@rm prybar
	@rm plugins/*