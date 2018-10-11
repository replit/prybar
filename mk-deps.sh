#/bin/bash

FOUND=""
check() {
	pkg-config --exists $2
	if [ $? -eq 0 ]; then
		echo "Found plugin $1 => $2" 1>&2
		FOUND="$FOUND plugins/$1.so"
	else
		echo "Couldnt find $1" 1>&2
	fi
}

for plugin in $(ls -1 src/languages/); do
	PKG=$(grep cgo.*pkg-config src/languages/$plugin/main.go | cut -d ':' -f 2)
	check $plugin $PKG
done


echo "plugins = $FOUND"


