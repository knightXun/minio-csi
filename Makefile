
.PHONY: all csi

CONTAINER_CMD?=docker

CSI_IMAGE_NAME=$(if $(ENV_CSI_IMAGE_NAME),$(ENV_CSI_IMAGE_NAME),quay.io/minio)
CSI_IMAGE_VERSION=$(if $(ENV_CSI_IMAGE_VERSION),$(ENV_CSI_IMAGE_VERSION),canary)

$(info csi image settings: $(CSI_IMAGE_NAME) version $(CSI_IMAGE_VERSION))

all: csi

test: go-test static-check

go-test:
	./scripts/test-go.sh

static-check:
	./scripts/lint-go.sh
	./scripts/lint-text.sh

.PHONY: csi
csi: ./build/build-oss.sh

image: csi
	cp _output/csi deploy/
	$(CONTAINER_CMD) build -t $(CSI_IMAGE_NAME):$(CSI_IMAGE_VERSION) deploy/

push-image: image
	$(CONTAINER_CMD) push $(CSI_IMAGE_NAME):$(CSI_IMAGE_VERSION)


clean:
	go clean -r -x
	rm -f deploy/csi/image/csi
	rm -f _output/csi
