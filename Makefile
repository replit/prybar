.PHONY: clean default

default: prybar

deps.lst: mk-deps.sh
	./mk-deps.sh > deps.lst

-include deps.lst

prybar: src/*.go $(plugins)
	cd src && go build -o ../prybar

plugins/%.so: src/languages/%/*.* 
	cd src/languages/$* && CGO_LDFLAGS_ALLOW=".*" go build -buildmode=plugin -o ../../../$@

clean:
	@rm deps.lst
	@rm prybar
	@rm plugins/*