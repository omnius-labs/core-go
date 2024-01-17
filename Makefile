install:
	make -C ./pkg/migration install

test:
	make -C ./pkg/base test
	make -C ./pkg/cloud_aws test
	make -C ./pkg/migration test
