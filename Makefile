.PHONY: clear-storage
clear-storage:
	rm -rf storage/bucket*/*

.PHONY: tests
tests:
	hurl --variable host=localhost --file-root testdata  --test tests/happy-path.hurl

.PHONY: test-after-add-new-node
test-after-add-new-node:
	hurl --variable host=localhost --test tests/after-add-new-bucket.hurl
