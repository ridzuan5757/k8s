URL := localhost:8080/station

# .PHONY: env
# env:
# 	@echo "Fetching data from $(URL)..."
# 	@curl -s $(URL) > raw_response.json
# 	@echo "Raw response:"
# 	@cat raw_response.json
# 	@echo "\nProcessing data..."
# 	@cat raw_response.json | sed 's/[{"}]//g; s/,/\n/g' > step1.txt
# 	@echo "After sed:"
# 	@cat step1.txt
# 	@cat step1.txt | \
# 		awk -F: '{ gsub(/^[ \t]+|[ \t]+$$/, "", $$1); \
# 		           gsub(/^[ \t]+|[ \t]+$$/, "", $$2); \
# 		           print $$1"="$$2 }' > .env
# 	@echo "\nContents of .env file:"
# 	@cat .env
# 	@rm raw_response.json step1.txt

.PHONY: env

env:
	@echo "Fetching data from $(URL)..."
	@curl -s $(URL) > raw_response.json
	@echo "Raw response:"
	@cat raw_response.json
	@echo "\nProcessing data..."
	@cat raw_response.json | sed 's/[{}]//g; s/,/\n/g' | \
		sed 's/^[[:space:]]*"//; s/"[[:space:]]*$$//; s/"[[:space:]]*:[[:space:]]*"/=/g; s/"[[:space:]]*:[[:space:]]*\([^"]\)/=\1/g' | \
		awk 'BEGIN {FS="="; in_outlet=0} \
		{ \
			gsub(/^[ \t]+|[ \t]+$$/, "", $$1); \
			gsub(/^[ \t]+|[ \t]+$$/, "", $$2); \
			if ($$1 == "outlet") { \
				in_outlet = 1; \
			} else if (in_outlet && $$1 ~ /^[a-zA-Z]/) { \
				print "outlet."$$1"="$$2; \
			} else if (!in_outlet) { \
				print $$0; \
			} \
		}' > .env
	@echo "\nContents of .env file:"
	@cat .env
	@rm raw_response.json



.PHONY: all
all: env
	docker-compose down && docker-compose up --build -d
