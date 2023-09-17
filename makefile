hello:
	echo "Hello"

build:
	@go get -d -v
	@go build -v .

# check for duplicate backends in a terraform file within the /Users/ directory
run-dev-tf:
	@INPUT_LOG_LEVEL=debug \
	GITHUB_ENV=./tmp/output.txt \
	GITHUB_OUTPUT=./tmp/output.txt \
	INPUT_CONTENT_REGEX="(?m)^.+backend\s+\"s3\"\s+{[^}]+key\s*=\s*\"([^\"]+)\"" \
	INPUT_DIRECTORY_SCOPE=/Users/ \
	INPUT_CHECK_FILE_EXTENSION=tf \
	go run -race .

# check for duplicate kafka consumer subscriptions to the same topic within the /Users/ directory
run-dev-kotlin:
	@INPUT_LOG_LEVEL=info \
	GITHUB_ENV=./tmp/output.txt \
	GITHUB_OUTPUT=./tmp/output.txt \
	INPUT_IGNORE_FILES="dev.conf,diff_test.conf" \
	INPUT_IGNORE_PATHS_CONTAINING="build/resources" \
	INPUT_CONTENT_REGEX="(?m)^kafka\.source\s*=\s*{\s*(.*\n)*\s*group_id\s*=\s*\"([^\"]*)\"" \
	INPUT_DIRECTORY_SCOPE=/Users/ \
	INPUT_CHECK_FILE_EXTENSION=conf \
	go run .
