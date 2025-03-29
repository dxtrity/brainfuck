$SRCFILE = Resolve-Path "./src/test/main.test.go"
$BUILD_DIR = Resolve-Path "./build"
$PNAME = Split-Path -Path $pwd -Leaf
$EXE = "$BUILD_DIR\$PNAME.test.exe"

# Ensure the build dir exists
if (-Not (Test-Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
}

# Build the Go binary
go build -o $EXE $SRCFILE

# Test if the build was succesful
if (Test-Path $EXE) {
    Write-Output "Build succesful."
} else {
    Write-Error "Build failed..."
}