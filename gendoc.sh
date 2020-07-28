#!/bin/bash
set -u

DOC_DIR=docs
PKG=github.com/E-Health/goscar-export

# Starting godoc server
echo "Starting godoc server..."
https_proxy='socks5://127.0.0.1:9090' godoc -http="${GO_DOC_HTTP:-:8085}" &
DOC_PID=$!

# Wait for the server to init (1minute, increase if reqd)
sleep 1m

# Scrape the pkg directory for the API docs. Scrap lib for the CSS/JS. Ignore everything else.
# The output is dumped to the directory "localhost:8085".
# wget -r -m -k -E -p -erobots=off --include-directories="/pkg/$PKG,/lib,/src/$PKG" --exclude-directories="*" "http://localhost:8085/pkg/$PKG/"
wget -r -m -k -E -p -erobots=off --include-directories="/pkg/$PKG,/lib" --exclude-directories="*" "http://localhost:8085/pkg/$PKG/"

# Stop the godoc server
kill -9 $DOC_PID

# Delete the old directory or else mv will put the localhost dir into
# the DOC_DIR if it already exists.
rm -rf $DOC_DIR
mv localhost\:8085 $DOC_DIR

echo "Docs can be found in $DOC_DIR"
echo "Replace /lib and /pkg in the gh-pages branch to update gh-pages"