GO_FOLDERS=./chiapi/... ./pkg/...

pre-commit: go-test go-lint go-dep-ensure go-doc

commit:
	@git cz

setup:
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@echo "#!/bin/bash\ncat \$$1 | commitlint" > .git/hooks/commit-msg
	@chmod +x .git/hooks/commit-msg
	@dep init

utilities:
	@go get -u github.com/robertkrimen/godocdown/godocdown

	@npm install -g eslint

	@npm install -g \
	commitizen \
	cz-conventional-changelog \
	@commitlint/prompt-cli \
	@commitlint/config-conventional

	@npm install -g \
	semantic-release \
	@semantic-release/commit-analyzer \
	@semantic-release/release-notes-generator \
	@semantic-release/changelog \
	@semantic-release/git

go-test:
	@go test $(GO_FOLDERS)

go-lint:
	@golangci-lint run \
	--exclude-use-default=false --disable-all \
	--enable golint --enable gosec --enable interfacer --enable unconvert \
	--enable goimports --enable goconst --enable gocyclo --enable misspell \
	--enable scopelint \
	$(GO_FOLDERS)

go-dep-ensure:
	@dep ensure

go-doc:
	@gfind -type d -printf '%d\t%P\n' | sort -r -nk1 | cut -f2- | \
		grep -v '^\.' | \
		grep -v '\/\.' | \
		grep -v '^pkg$$' | \
		grep -v '^vendor' | \
		xargs -I{} bash -c "godocdown {} > {}/README.md"

guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi
