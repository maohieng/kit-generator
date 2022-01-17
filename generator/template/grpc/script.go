package grpc

import "fmt"

const (
	scriptTemplate = `#!/usr/bin/env sh

# Install proto3
# sudo apt-get install -y git autoconf automake libtool curl make g++ unzip
# git clone https://github.com/google/protobuf.git
# cd protobuf/
# ./autogen.sh
# ./configure
# make
# make check
# sudo make install
# sudo ldconfig # refresh shared library cache.
#
# Update protoc Go bindings via
#  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#
# See also
#  https://github.com/grpc/grpc-go/tree/master/examples

protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative %s`

	windowsScriptTemplate = `:: Install proto3.
:: https://github.com/google/protobuf/releases
:: Update protoc Go bindings via
::  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
::  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
::
:: See also
::  https://github.com/grpc/grpc-go/tree/master/examples

protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative %s`

	darwinScriptText = `#!/usr/bin/env sh

# Install proto3 from source macOS only.
#  brew install autoconf automake libtool
#  git clone https://github.com/google/protobuf
#  ./autogen.sh ; ./configure ; make ; make install
#
# Update protoc Go bindings via
#  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# See also
#  https://github.com/grpc/grpc-go/tree/master/examples

protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative %s`
)

// WindowsScriptText return gRPC compile script text for windows,
// it uses the syntax of windows batch
func WindowsScriptText(name string) string {
	return fmt.Sprintf(windowsScriptTemplate, name)
}

// DarwinScriptText return gRPC compile script text for macOS, it
// uses the syntax of shell script
func DarwinScriptText(name string) string {
	return fmt.Sprintf(darwinScriptText, name)
}

// ScriptText return gRPC compile script text, it uses the syntax
// of shell script
func ScriptText(name string) string {
	return fmt.Sprintf(scriptTemplate, name)
}
