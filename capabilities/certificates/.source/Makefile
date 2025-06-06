VERSION ?= latest

#
# code generation
#
WORKLOAD_CONFIG ?= .source/workload-cert-manager.yaml
define operator_builder_init
	operator-builder init \
		--workload-config $(WORKLOAD_CONFIG) \
		--repo github.com/klearwave/operators-capability/capabilities/certificates \
		--controller-image ghcr.io/klearwave/certificates-operator:$(VERSION) \
		--skip-go-version-check
endef

define operator_builder_create
	operator-builder create api \
		--workload-config $1 \
		--controller \
		--resource \
		--force
endef

MANAGED_ASSETS=../go.mod \
	../README.md \
	../.gitignore \
	../.dockerignore \
	../.github/workflows/release.yml \
	../.goreleaser.yml \
	../Dockerfile \
	../LICENSE \
	../config/default/kustomization.yaml \
	../config/default/manager_auth_proxy_patch.yaml \
	../controllers/certificates/certmanager_phases.go \
	../controllers/certificates/trustmanager_phases.go \
	../apis/certificates/v1alpha1/certmanager/mutate/mutating_webhook_cert_manager_webhook.go \
	../apis/certificates/v1alpha1/certmanager/mutate/validating_webhook_cert_manager_webhook.go \
	../apis/certificates/v1alpha1/trustmanager/mutate/validating_webhook_trust_manager.go

# preserve manually mananged assets
preserve:
	@mkdir -p assets; \
	for ASSET in $(MANAGED_ASSETS); do \
		DEST=`echo $$ASSET | sed 's/\.\.\///g'`; \
		if [ -f $$ASSET ]; then \
			echo "preserving file $$ASSET"; \
			DIR=`dirname $$DEST`; \
			mkdir -p assets/$$DIR; \
			cp -R $$ASSET assets/$$DEST; \
		elif [ -d $$ASSET ]; then \
			echo "preserving directory $$ASSET"; \
			mkdir -p assets/$$DEST; \
			cp -R $$ASSET/* assets/$$DEST; \
		else \
			echo "not preserving non-existing asset $$ASSET"; \
		fi \
	done

# restore manually managed assets
restore:
	@cp -R assets/* ..

operator-clean: preserve
	sudo rm -rf ../bin && rm -rf ../*

operator-init:
	cd .. && $(call operator_builder_init)

operator-create:
	cd .. && \
		for workload in $$(find .source -maxdepth 1 -type f -name "workload-*.yaml"); do \
			$(call operator_builder_create, "$${workload}") ;\
		done && \
	make manifests

operator: download-all overlay-all operator-clean operator-init operator-create restore controller-manifests
	cd ../ && go mod tidy

#
# project management
#

# create scaffold
PROJECT ?=
project:
	@mkdir -p $(PROJECT)/config $(PROJECT)/vendor $(PROJECT)/manifests
	@touch $(PROJECT)/config/values.yaml $(PROJECT)/config/vendor.yaml $(PROJECT)/config/overlays.yaml

# download from upstream
download = (vendir sync --file $1/config/vendor.yaml --lock-file $1/config/vendor.yaml.lock)
download:
	$(call download, $(PROJECT))

download-all:
	@for project in $$(find . -maxdepth 1 -type d ! -name "assets" ! -name "." -exec basename {} \;); do \
		$(call download, $$project); \
	done

# overlay changes
overlay = (yot -I 2 -i $1/config/overlays.yaml -o $1 -f $1/config/values.yaml)
overlay:
	$(call overlay, $(PROJECT))

overlay-all:
	@for project in $$(find . -maxdepth 1 -type d ! -name "assets" ! -name "." -exec basename {} \;); do \
		$(call overlay, $$project); \
	done

# controller manifests
# NOTE: we need to preserve the original kustomization file as we always want our generated manifest to 
#       be set to the latest release except when we tag a new version.
CONTROLLER_IMG ?= ghcr.io/klearwave/certificates-operator:latest
MANIFESTS_FILE ?= ../../config/deploy/manifests.yaml
controller-manifests:
	cd .. && \
	make kustomize && \
	mkdir -p config/deploy && \
	cd config/manager && \
	cp -f kustomization.yaml kustomization.yaml.bak && \
	export KUSTOMIZE_CMD=$$(ls ../../bin/kustomize*) && \
	$$KUSTOMIZE_CMD edit set image controller=$(CONTROLLER_IMG) && \
	$$KUSTOMIZE_CMD build ../../config/default > $(MANIFESTS_FILE) && \
	mv kustomization.yaml.bak kustomization.yaml
