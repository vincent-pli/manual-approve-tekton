module github.com/vincentpli/manual-approve-tekton

go 1.15

require (
	github.com/tektoncd/pipeline v0.30.0
	go.uber.org/zap v1.19.1
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	k8s.io/code-generator v0.21.4
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7
	knative.dev/hack v0.0.0-20211122162614-813559cefdda
	knative.dev/hack/schema v0.0.0-20211122162614-813559cefdda
	knative.dev/pkg v0.0.0-20211125172117-608fc877e946
	knative.dev/sample-controller v0.27.0
)
