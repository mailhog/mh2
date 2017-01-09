VERSION := 2.0.0-alpha
BUILD_DATE := `date +%Y-%m-%d\ %H:%M`
VERSIONFILE := version/version.go

all: version test
	go install ./cmd/...

test:
	go test ./...

version:
	# FIXME -ldflags would be better to avoid source changes on build
	# http://grokbase.com/t/gg/golang-nuts/14c4dtb7ta/go-nuts-using-ldflags-to-set-variables-in-package-other-than-main
	# https://stackoverflow.com/questions/11354518/golang-application-auto-build-versioning
	rm -f $(VERSIONFILE)
	@echo "package version" > $(VERSIONFILE)
	@echo "import \"fmt\"" >> $(VERSIONFILE)
	@echo "const (" >> $(VERSIONFILE)
	@echo "  // Version is the version" >> $(VERSIONFILE)
	@echo "  Version = \"$(VERSION)\"" >> $(VERSIONFILE)
	@echo "  // BuildDate is the build date" >> $(VERSIONFILE)
	@echo "  BuildDate = \"$(BUILD_DATE)\"" >> $(VERSIONFILE)
	@echo ")" >> $(VERSIONFILE)
	@echo "// String returns a formatted version and build date string" >> $(VERSIONFILE)
	@echo "func String() string {" >> $(VERSIONFILE)
	@echo " return fmt.Sprintf(\"%s (%s)\", Version, BuildDate)" >> $(VERSIONFILE)
	@echo "}" >> $(VERSIONFILE)
	go fmt $(VERSIONFILE)

debug: all
	HUMAN_LOG=1 mh2

.PHONY: all test version debug
