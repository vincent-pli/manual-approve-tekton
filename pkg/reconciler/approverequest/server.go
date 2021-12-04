/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package approverequest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	cache "github.com/patrickmn/go-cache"
	approverequestsv1alpha1 "github.com/vincent-pli/manual-approve-tekton/pkg/apis/approverequests/v1alpha1"
	approverequestclientset "github.com/vincent-pli/manual-approve-tekton/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RequestApprove struct {
	ApproveTemplate string
	Requests        []RequestPlain
}

type RequestPlain struct {
	RequestName      string
	RequestTimestamp string
	Approved         bool
	ApproveTimestamp string
}

type WebServer struct {
	cache                   *cache.Cache
	approverequestClientSet approverequestclientset.Interface
}

func (ws *WebServer) Start(context context.Context) error {
	http.HandleFunc("/requsts", ws.listRequest)
	http.HandleFunc("/approve", ws.approveRequest)
	http.ListenAndServe(":8090", nil)
	return nil
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (ws *WebServer) approveRequest(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	approvetemplates, ok := req.URL.Query()["approvetemplate"]
	if !ok || len(approvetemplates[0]) < 1 {
		log.Println("Url Param 'approvetemplate' is missing")
		return
	}

	requests, ok := req.URL.Query()["request"]
	if !ok || len(requests[0]) < 1 {
		log.Println("Url Param 'request' is missing")
		return
	}

	namespacename := strings.Split(approvetemplates[0], "/")

	ar, err := ws.approverequestClientSet.CustomV1alpha1().ApproveRequests(namespacename[0]).Get(context.TODO(), namespacename[1], metav1.GetOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	arCopy := ar.DeepCopy()
	for index, req := range arCopy.Status.Requests {
		if req.RequestName == requests[0] {
			arCopy.Status.Requests[index].Approved = true
		}
	}

	_, err = ws.approverequestClientSet.CustomV1alpha1().ApproveRequests(namespacename[0]).UpdateStatus(context.TODO(), arCopy, metav1.UpdateOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "done\n")
}

func (ws *WebServer) listRequest(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	requestApproves := []RequestApprove{}
	for key, item := range ws.cache.Items() {
		requestApprove := RequestApprove{}
		requestApprove.ApproveTemplate = key
		reqs, ok := item.Object.([]approverequestsv1alpha1.Request)
		if ok {
			requestApprove.Requests = translate(reqs)
		}
		requestApproves = append(requestApproves, requestApprove)
	}
	payload, err := json.Marshal(requestApproves)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Fprintf(w, "hello\n")
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func translate(orgReqs []approverequestsv1alpha1.Request) []RequestPlain {
	requests := []RequestPlain{}

	for _, req := range orgReqs {
		request := RequestPlain{}
		request.RequestName = req.RequestName
		request.RequestTimestamp = req.RequestTimestamp.String()
		request.Approved = req.Approved
		request.ApproveTimestamp = req.ApproveTimestamp.String()
		requests = append(requests, request)
	}

	return requests
}
