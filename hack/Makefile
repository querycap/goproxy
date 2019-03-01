DOCKERX_CONTEXT ?= .
DOCKERX_NAME ?= default
DOCKERX_OUTPUT ?=
DOCKERX_PUSH ?= false
DOCKERX_ARCH_SUFFIX ?= false
DOCKERX_PLATFORMS ?= linux/amd64 linux/arm64
DOCKERX_BUILD_ARGS ?=
DOCKERX_LABELS ?=
DOCKERX_TAGS ?= latest
DOCKERX_TAG_SUFFIX ?=

ifeq ($(DOCKERX_PUSH),true)
	DOCKERX_OUTPUT = --push
endif

dockerx:
	@set -eux; \
	\
	docker buildx build $(DOCKERX_OUTPUT) \
		$(foreach h,$(HUB),$(foreach t,$(DOCKERX_TAGS),--tag=$(h)/$(DOCKERX_NAME):$(t)$(DOCKERX_TAG_SUFFIX))) \
		$(foreach p,$(DOCKERX_PLATFORMS),--platform=$(p)) \
		$(foreach a,$(DOCKERX_BUILD_ARGS),--build-arg=$(a)) \
		$(foreach l,$(DOCKERX_LABELS),--label=$(l)) \
		--file $(DOCKERX_CONTEXT)/Dockerfile.$(DOCKERX_NAME) $(DOCKERX_CONTEXT)

