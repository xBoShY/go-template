OAPI := ./generated/model/model.yaml
OAPI += ./generated/routes/routes.yaml
OAPI += ./generated/client/client.yaml

default: build

build:
	@for dir in $(OAPI); do \
		echo "$$dir"; \
		oapi-codegen -config $$dir oapi.yaml; \
	done
