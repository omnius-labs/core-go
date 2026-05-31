modules = \
	aws \
	base \
	migration \
	testing

install:
	for dir in $(modules); do \
        make -C ./$$dir install; \
    done

test:
	for dir in $(modules); do \
        make -C ./$$dir test; \
    done
