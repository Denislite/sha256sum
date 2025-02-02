package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AdmissionReviewFromRequest(r *http.Request) (*admissionv1.AdmissionReview, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("expected application/json content-type")
	}

	admissionReviewRequest := &admissionv1.AdmissionReview{}

	err := json.NewDecoder(r.Body).Decode(&admissionReviewRequest)
	if err != nil {
		return nil, err
	}
	return admissionReviewRequest, nil
}

func AdmissionResponseFromReview(admReview *admissionv1.AdmissionReview) (*admissionv1.AdmissionResponse, error) {
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admReview.Request.Resource != podResource {
		err := fmt.Errorf("did not receive pod, got %s", admReview.Request.Resource.Resource)
		return nil, err
	}

	admissionResponse := &admissionv1.AdmissionResponse{}

	rawRequest := admReview.Request.Object.Raw
	pod := corev1.Pod{}

	err := json.NewDecoder(bytes.NewReader(rawRequest)).Decode(&pod)
	if err != nil {
		err := fmt.Errorf("error decoding raw pod: %v", err)
		return nil, err
	}

	var patch string
	patchType := v1.PatchTypeJSONPatch

	log.Println(pod.Name)

	log.Println("pod has following labels", pod.Labels)
	if _, ok := pod.Labels["tcpdump-sidecar"]; ok {
		patch = `[
		{
			"op":"add",
			"path":"/spec/containers/1",
			"value":{
				"image":"hasher:latest",
				"imagePullPolicy":"Never",
				"name":"hasher",
				"envFrom": [
				  {
					"configMapRef": {
					  "name": "hasher-config"
					}
				  },
				  {
					"secretRef": {
					  "name": "postgres-secret"
					}
				  }
				],
				"command": [
				  "sh",
				  "-c",
                  "conf=$(ls /etc/hasher-config); data=$(cat /etc/hasher-config/$conf);env1=$(echo $data | cut -f 1 -d\" \"); PID_NAME=\"${env1#*=}\";env2=$(echo $data | cut -f 2 -d\" \"); MOUNT_PATH=\"${env2##*=}\";pid=$(pidof -s $PID_NAME); ./sha256sum -d ../proc/$pid/root/$MOUNT_PATH;"
				],
				"volumeMounts": [
                  {
      				"name": "hasher-config",
  					"mountPath": "/etc/hasher-config",
      				"readOnly": true
   				  }
  				],
				"env": [
				  {
					"name": "POD_NAME",
					"valueFrom": {
					  "fieldRef": {
						"fieldPath": "metadata.name"
					  }
					}
				  }
				],
				"resources": {
				"limits": {
				  "memory": "50Mi",
				  "cpu": "50m"
				}
			  },
			  "securityContext": {
				"capabilities": {
				  "add": [
					"SYS_PTRACE"
				  ]
				}
			  },
			  "stdin": true,
			  "tty": true
			}
		}
	]`
	}

	admissionResponse.Allowed = true
	if patch != "" {
		log.Println("patching the pod with:", patch)
		admissionResponse.PatchType = &patchType
		admissionResponse.Patch = []byte(patch)
	}

	return admissionResponse, nil
}
