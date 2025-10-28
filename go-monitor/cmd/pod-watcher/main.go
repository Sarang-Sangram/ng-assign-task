package main

import (
    "fmt"
    "time"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/cache"
    "os"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/cache"
)

func main() {
    var config *rest.Config
    var err error
    if kube := os.Getenv("KUBECONFIG"); kube != "" {
        config, err = clientcmd.BuildConfigFromFlags("", kube)
        if err != nil {
            panic(err)
        }
    } else {
        config, err = rest.InClusterConfig()
        if err != nil {
            panic(err)
        }
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    factory := informers.NewSharedInformerFactory(clientset, 0)
    podInformer := factory.Core().V1().Pods().Informer()

    podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            p := obj.(*corev1.Pod)
            fmt.Printf("[%s] Pod created: %s/%s\n", time.Now().Format(time.RFC3339), p.Namespace, p.Name)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            p := newObj.(*corev1.Pod)
            fmt.Printf("[%s] Pod updated: %s/%s phase=%s\n", time.Now().Format(time.RFC3339), p.Namespace, p.Name, p.Status.Phase)
        },
        DeleteFunc: func(obj interface{}) {
            var p *corev1.Pod
            switch t := obj.(type) {
            case *corev1.Pod:
                p = t
            case cache.DeletedFinalStateUnknown:
                if pod, ok := t.Obj.(*corev1.Pod); ok {
                    p = pod
                }
            }
            if p != nil {
                fmt.Printf("[%s] Pod deleted: %s/%s\n", time.Now().Format(time.RFC3339), p.Namespace, p.Name)
            }
        },
    })

    stopCh := make(chan struct{})
    factory.Start(stopCh)
    factory.WaitForCacheSync(stopCh)
    select {}
}
