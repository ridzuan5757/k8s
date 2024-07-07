# Determine the operating system
ifeq ($(OS),Windows_NT)
    DETECTED_OS := windows
else
    DETECTED_OS := $(shell uname -s)
    ifeq ($(DETECTED_OS),Linux)
        DETECTED_OS := linux
    else ifeq ($(DETECTED_OS),Darwin)
        DETECTED_OS := macos
    endif
endif

# Default target
all: check_os

# Target to check and display the OS
check_os:
	@echo "Detected Operating System: $(DETECTED_OS)"

# OS-specific targets
windows:
	@echo "Running Windows-specific commands"

linux:
	@echo "Running Linux-specific commands"

macos:
	@echo "Running macOS-specific commands"

# Target to run OS-specific commands
run_os_specific: 
	check_os
ifeq ($(DETECTED_OS),Windows)
	@$(MAKE) windows
else ifeq ($(DETECTED_OS),Linux)
	@$(MAKE) linux
else ifeq ($(DETECTED_OS),macOS)
	@$(MAKE) macos
else
	@echo "Unsupported operating system"
endif
