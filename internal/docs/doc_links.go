package docs

var FieldDocLinks = map[string]string{
	"Pod.spec.containers":                   "https://kubernetes.io/docs/concepts/workloads/pods/#pod-templates",
	"Pod.spec.volumes":                      "https://kubernetes.io/docs/concepts/storage/volumes/",
	"Pod.spec.containers[].livenessProbe":   "https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/",
	"Pod.spec.containers[].readinessProbe":  "https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/",
	"Pod.spec.containers[].resources":       "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
	"Pod.spec.containers[].imagePullPolicy": "https://kubernetes.io/docs/concepts/containers/images/#imagepullpolicy",

	"Deployment.spec.replicas":                                                "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#scaling-a-deployment",
	"Deployment.spec.strategy":                                                "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy",
	"Deployment.spec.strategy.rollingUpdate.maxUnavailable":                   "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#rolling-update-deployment",
	"Deployment.spec.selector":                                                "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#label-selector",
	"Deployment.spec.template":                                                "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#pod-template",
	"Deployment.spec.template.spec.containers[].securityContext.runAsNonRoot": "https://kubernetes.io/docs/tasks/configure-pod-container/security-context/",
	"Deployment.spec.template.spec.affinity":                                  "https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity",
	"Deployment.spec.template.spec.tolerations":                               "https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/",
	"Deployment.spec.template.spec.nodeSelector":                              "https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector",
	"Deployment.spec.template.metadata.annotations":                           "https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#kubectl-kubernetes-io-restartedat",

	"Service.spec.ports":    "https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service",
	"Service.spec.selector": "https://kubernetes.io/docs/concepts/services-networking/service/#services-without-selectors",
	"Service.spec.type":     "https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types",

	"Ingress.spec.rules": "https://kubernetes.io/docs/concepts/services-networking/ingress/#ingress-rules",
	"Ingress.spec.tls":   "https://kubernetes.io/docs/concepts/services-networking/ingress/#tls",

	"ConfigMap.data": "https://kubernetes.io/docs/concepts/configuration/configmap/#configmap-object",
	"Secret.data":    "https://kubernetes.io/docs/concepts/configuration/secret/#overview-of-secrets",

	"PersistentVolumeClaim.spec.resources": "https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims",
	"Namespace.metadata.name":              "https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/#working-with-namespaces",
}
