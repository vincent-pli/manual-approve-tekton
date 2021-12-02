module github.com/vincent-pli/manual-approve-tekton

go 1.15

require (
	github.com/hashicorp/go-multierror v1.1.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/tektoncd/pipeline v0.30.0
	go.uber.org/zap v1.19.1
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	k8s.io/code-generator v0.21.4
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7
	knative.dev/hack v0.0.0-20211122162614-813559cefdda
	knative.dev/hack/schema v0.0.0-20211122162614-813559cefdda
	knative.dev/pkg v0.0.0-20211125172117-608fc877e946
)
