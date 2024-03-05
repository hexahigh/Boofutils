all: main nms

main:
	go build -o boofutils main.go

nms:
	@ORIG_DIR=$$(pwd)
	@echo "Cloning repository..."

	@if [ ! -d "/tmp/libnms-build" ]; then \
		git clone https://github.com/bartobri/libnms.git "/tmp/libnms-build"; \
	else \
		echo "Directory /tmp/libnms-build already exists (Somehow)"; \
	fi
	
	@echo "Building NMS..."
	@cd "/tmp/libnms" && make -f /tmp/libnms-build/Makefile && sudo make install -f /tmp/libnms-build/Makefile
	@rm -rf "/tmp/libnms-build"
	@echo "Done!"

check_tools:
	@echo "Checking for required tools to build nms..."
	@if ! command -v make &> /dev/null; then echo "make not found"; exit 1; fi
	@if ! command -v gcc &> /dev/null; then echo "gcc not found"; exit 1; fi
	@if ! command -v git &> /dev/null; then echo "git not found"; exit 1; fi
	@echo "All required tools found."