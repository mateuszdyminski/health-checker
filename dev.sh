#!/bin/bash

if [ -f .app ]; then
	source ./.app
fi

usage() {
	cat <<EOF
Usage: $(basename $0) <command>

Wrappers around core binaries:
    app                    Runs the web app in development mode.
    deps                   Installs dependencies.
    release                Builds release package.

EOF
	exit 1
}

GO=${GO:-$(which go)}
export GOPATH="$PWD"
export PATH="$PATH:$GOPATH/bin"

# Absolute URL configuration
export EXT_STATIC_DIR="$PWD/static"

BUILDDIR="$PWD/build"
LOGDIR=/tmp/app/logs

install_deps() {
	echo "Installing dependencies..."
	$GO get github.com/gorilla/websocket
	$GO get github.com/gorilla/mux
	$GO get github.com/golang/glog
}

release() {
	set -e
	install_deps
	export GOOS=linux
	export GOARCH=arm
	BUILD_OUT="$BUILDDIR/app"
	echo "Building for $GOOS/$GOARCH in $BUILD_OUT"
	rm -rf "$BUILD_OUT"
	mkdir -p "$BUILD_OUT/db" "$BUILD_OUT/static"
	$GO build -o "$BUILD_OUT/app" app/app
	cp dev.sh "$BUILD_OUT/"
	cp -r app/* "$BUILD_OUT/static/"
	set +e
	cd "$(dirname "$BUILD_OUT")" && tar zcf "$BUILD_OUT.tgz" "$(basename "$BUILD_OUT")"
	echo "DONE. $BUILD_OUT.tgz"
}

CMD="$1"
shift
case "$CMD" in
	deps)
		install_deps
	;;
	app)
		build
		$GO build -o "bin/app" app/app
		flags=$(test -z $EXT_STATIC_DIR)
		set -x
		exec "$GOPATH/bin/app" \
			--log_dir="$LOGDIR" \
			--alsologtostderr \
			--stderrthreshold=INFO \
			--dir="$EXT_STATIC_DIR" \
			--host="localhost" \
			--port="8090" \
			--address="https://bitbucket.org/Avaus/kone-china-web/pull-requests" \
			$flags \
			$@ \
	;;
	release)
		release
	;;
	build)
		set -e
		install_deps
		$GO build -o "bin/app" app/app
		set +e
		echo "Build success, bin directory: bin/app"
	;;
	*)
		usage
	;;
esac