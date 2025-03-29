$SRCFILE = Resolve-Path "./src/main.go"
$BUILD_DIR = Resolve-Path "./build"
$EXE = "$BUILD_DIR\super.exe"

Write-Output "Source File: $SRCFILE"
Write-Output "Build Directory: $BUILD_DIR"
Write-Output "Executable Path: $EXE"

# Ensure the build dir exists
if (-Not (Test-Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
}

# Build the Go binary
go build -o $EXE $SRCFILE

# Run the Go binary
if (Test-Path $EXE) {
    Write-Output "Build succesful. Running executable..."
    & $EXE
} else {
    Write-Error "Build failed..."
}