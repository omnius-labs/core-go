modules = \
	aws \
	base \
	migration \
	testing

install:
	go install golang.org/x/vuln/cmd/govulncheck@0782b76014f15f24e22a438f30f308df42899ba1
	go install github.com/onsi/ginkgo/v2/ginkgo@4f62d7a74752034222d97d911f904d9be47ff7aa

test:
	for dir in $(modules); do \
		(cd ./$$dir && ginkgo -r) || exit $$?; \
    done

audit:
	for dir in $(modules); do \
        (cd ./$$dir && govulncheck ./...) || exit $$?; \
    done
