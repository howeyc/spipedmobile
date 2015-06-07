#!/bin/bash -e

# Build the go library
ORIG=$(pwd)

# Check for GOLANG installation
if [ -z $GOROOT ] || [[ $(go version) != go\ version\ go1.4beta1* ]] ; then
        mkdir -p "golib"
        tmpgo='golib/go'
        if [ ! -f "$tmpgo/bin/go" ]; then
                # Download GOLANG v1.3.3
                wget -O go.src.tar.gz https://golang.org/dl/go1.4rc1.src.tar.gz
                sha1=$(sha1sum go.src.tar.gz)
                if [ "$sha1" != "ff8e7d78e85658251a36e45f944af70f226368ab  go.src.tar.gz" ]; then
                        echo "go.src.tar.gz SHA1 checksum does not match!"
                        exit 1
                fi
                mkdir -p $tmpgo
                tar -xzf go.src.tar.gz --strip=1 -C $tmpgo
                rm go.src.tar.gz
                # Build GO for host
                pushd $tmpgo/src
                ./make.bash --no-clean
                popd
        fi
        # Add GO to the environment
        export GOROOT="$(pwd)/$tmpgo"
fi

# Make toolchain
wget http://dl.google.com/android/ndk/android-ndk-r10e-linux-x86.bin
chmod a+x android-ndk-r10e-linux-x86.bin
$ORIG/android-ndk-r10e-linux-x86.bin

mkdir -p golib/ndk-toolchain
$ORIG/android-ndk-r10e/build/tools/make-standalone-toolchain.sh --platform=android-9 --install-dir=golib/ndk-toolchain --toolchain=arm-linux-androideabi-4.8

# Check whether GOLANG is compiled with cross-compilation for arm
if [ ! -f $GOROOT/bin/android_arm/go ]; then
        pushd $GOROOT/src
        # Build GO for cross-compilation
        CC_FOR_TARGET=$ORIG/golib/ndk-toolchain/bin/arm-linux-androideabi-gcc GOOS=android GOARCH=arm ./make.bash --no-clean
        popd
fi

# Setup GOPATH
mkdir -p "golib/gopath"
cd "golib/gopath"
export GOPATH="$(pwd)"

# Get needed build tools
$GOROOT/bin/go get github.com/alecthomas/gobundle/gobundle
$GOROOT/bin/go get golang.org/x/mobile/cmd/gobind

# Get library (Will fail, need to compile assets)
set +e
$GOROOT/bin/go get github.com/howeyc/spipedmobile
set -e

# Setup library
pushd $GOPATH/src/github.com/howeyc/spipedmobile
pushd webroot
$GOPATH/bin/gobundle --recursive --compress --uncompress_on_init --retain_uncompressed --bundle="webroot" --package=assets --target=../assets/bundle.go bootstrap-3.3.0 jquery.min.js license.html template.index.html
popd
mkdir -p $ORIG/src/main/java/go/spiped
CGO_ENABLED=0 $GOPATH/bin/gobind -lang=java github.com/howeyc/spipedmobile > $ORIG/src/main/java/go/spiped/Spipedmobile.java
mkdir -p go_spiped
CGO_ENABLED=0 $GOPATH/bin/gobind -lang=go github.com/howeyc/spipedmobile > go_spiped/go_spiped.go
pushd android_so_lib
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 $GOROOT/bin/go build -o $ORIG/src/main/jniLibs/armeabi/libgojni.so -ldflags="-shared"
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=7 $GOROOT/bin/go build -o $ORIG/src/main/jniLibs/armeabi-v7a/libgojni.so -ldflags="-shared"
popd
popd

cd $ORIG
rm -rf golib
