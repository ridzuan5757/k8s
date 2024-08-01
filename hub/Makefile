URL := localhost:8080/station

.PHONY: env
env:
	@echo "Fetching data from $(URL)..."
	@curl -s $(URL) > raw_response.json
	@echo "Raw response:"
	@cat raw_response.json
	@echo "\nProcessing data..."
	@cat raw_response.json | sed 's/[{"}]//g; s/,/\n/g' > step1.txt
	@echo "After sed:"
	@cat step1.txt
	@cat step1.txt | \
		awk -F: '{ gsub(/^[ \t]+|[ \t]+$$/, "", $$1); \
		           gsub(/^[ \t]+|[ \t]+$$/, "", $$2); \
		           print $$1"="$$2 }' > .env
	@echo "\nContents of .env file:"
	@cat .env
	@rm raw_response.json step1.txt

.PHONY: test
test:
	@echo "Using variables from .env file:"
	@grep -v '^#' .env | xargs -I {} echo "Variable: {}"

.PHONY: all
all: env test
	docker-compose down && docker-compose up --build -d
