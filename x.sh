#/bin/bash
for item in net sys crypto oauth2 tools perf mobile blog lint vgo time text review image term sync exp talks; do
	git clone https://github.com/golang/${item} $GOPATH/src/golang.org/x/${item}
done
