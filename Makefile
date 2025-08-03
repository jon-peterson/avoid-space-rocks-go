.PHONY: lint test build clean build-windows build-macos build-linux all-platforms

# Default target
all: lint build

# Lint the code
lint:
	@staticcheck ./...

# Run tests
test:
	@go test ./...

# Clean build artifacts
clean:
	@rm -rf bin/
	@find . -name "*.test" -delete
	@find . -name "*.out" -delete

# Clean downloaded files
clean-downloads:
	@rm -rf bin/downloads/

# Build for current platform
build: test
	@mkdir -p bin
	@cd cmd/avoid_space_rocks && go build -o ../../bin/avoid-space-rocks .

# Cross-compilation targets

# Download raylib Windows DLL if it doesn't exist
raylib-windows:
	@mkdir -p bin/downloads/raylib-windows
	@if [ ! -f bin/downloads/raylib-windows/raylib-5.5_win64_mingw-w64/lib/raylib.dll ]; then \
		echo "Downloading raylib Windows DLL..."; \
		curl -L -o bin/downloads/raylib-5.5.0-w64-mingw.zip https://github.com/raysan5/raylib/releases/download/5.5/raylib-5.5_win64_mingw-w64.zip; \
		unzip -o bin/downloads/raylib-5.5.0-w64-mingw.zip -d bin/downloads/raylib-windows; \
	fi

# Build for Windows (64-bit)
build-windows: test raylib-windows
	@mkdir -p bin/windows
	@cd cmd/avoid_space_rocks && \
		GOOS=windows GOARCH=amd64 go build -o ../../bin/windows/avoid-space-rocks.exe .

# Download raylib Linux libraries if they don't exist
raylib-linux:
	@mkdir -p bin/downloads/raylib-linux
	@if [ ! -f bin/downloads/raylib-linux/raylib-5.5_linux_amd64/lib/libraylib.so ]; then \
		echo "Downloading raylib Linux libraries..."; \
		curl -L -o bin/downloads/raylib-linux/raylib-5.5.0_linux_amd64_from_build.tar.gz https://github.com/raysan5/raylib/releases/download/5.5/raylib-5.5_linux_amd64.tar.gz; \
		tar -xzf bin/downloads/raylib-linux/raylib-5.5.0_linux_amd64_from_build.tar.gz -C bin/downloads/raylib-linux/; \
	fi

# Build for Linux (64-bit)
build-linux: test raylib-linux
	@mkdir -p bin/linux
	@cd cmd/avoid_space_rocks && \
		CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
		go build -ldflags="-s -w" \
		-o ../../bin/linux/avoid-space-rocks .

# Build for all platforms
all-platforms: build-windows build-linux

# Package release (creates zip files with flat structure)
package: all-platforms
	@echo "Creating release packages..."
	@# Create temporary directories with flat structure
	@mkdir -p bin/package-linux
	@mkdir -p bin/package-windows

	@# Copy Windows files
	@cp bin/windows/avoid-space-rocks.exe bin/package-windows/
	@cp bin/downloads/raylib-windows/raylib-5.5_win64_mingw-w64/lib/raylib.dll bin/package-windows/

	@# Copy Linux files and libraries
	@cp bin/linux/avoid-space-rocks bin/package-linux/
	@cp bin/downloads/raylib-linux/raylib-5.5_linux_amd64/lib/libraylib.so* bin/package-linux/

	@# Copy assets to both packages
	@cp -r assets bin/package-linux/
	@cp -r assets bin/package-windows/

	@# Create zip files
	@cd bin/package-linux && zip -r ../avoid-space-rocks-linux.zip .
	@cd bin/package-windows && zip -r ../avoid-space-rocks-windows.zip .

	@# Clean up temporary directories
	@rm -rf bin/package-linux bin/package-windows
	@echo "Packages created successfully"

# List available targets
help:
	@echo "Available targets:"
	@echo "  build        - Build for current platform"
	@echo "  build-windows - Build for Windows (64-bit)"
	@echo "  build-linux   - Build for Linux (64-bit)"
	@echo "  all-platforms - Build for all platforms"
	@echo "  package      - Create release packages"
	@echo "  test         - Run tests"
	@echo "  lint         - Run static analysis"
	@echo "  clean        - Remove build artifacts"

