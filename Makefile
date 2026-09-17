.PHONY: proto
proto:
	protoc --go_out=$(OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_OUT) --go-grpc_opt=paths=source_relative \
		$(DIR_PROTO)

.PHONY: proto-deps
proto-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpcD@v1.2

.PHONY: gen
gen:
	protoc \
      --go_out=./pkg/protos/notification --go_opt=paths=source_relative \
      --go-grpc_out=./pkg/protos/notification --go-grpc_opt=paths=source_relative \
      --proto_path=./pkg/protos \
      ./pkg/protos/notification.proto
# ---------- Local k8s (minikube) ----------
MINIKUBE ?= minikube
KUBECTL  ?= kubectl
NS       := talker
IMG_TAG  := dev
SERVICES := user message chat

.PHONY: cluster-up cluster-down images load deploy undeploy status logs

cluster-up:
	$(MINIKUBE) start --driver=docker

cluster-down:
	$(MINIKUBE) stop

images:
	@for s in $(SERVICES); do \
		docker build -f services/$$s/Dockerfile -t talker/$$s:$(IMG_TAG) . ; \
	done

load:
	@for s in $(SERVICES); do \
		$(MINIKUBE) image load talker/$$s:$(IMG_TAG) ; \
	done

deploy:
	$(KUBECTL) apply -f deployment/

undeploy:
	$(KUBECTL) delete -f deployment/ --ignore-not-found

status:
	$(KUBECTL) -n $(NS) get pods

logs:
	$(KUBECTL) -n $(NS) logs -l app=chat-service --tail=50

# up via docker
docker-up:
	@for s in $(SERVICES); do \
		docker-compose -f services/$$s/docker-compose.yml up -d; \
	done

# up via k8s
k8s-up:
	@for s in $(SERVICES); do \
		kubectl apply -f deployment/$$s; \
	done
