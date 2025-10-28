We have below steps for AKS and helmcharts separately along with code for podwatcher

Deploy K8S Cluster:

1.	Deployed an AKS cluster with 2 nodes in Azure via portal. 
2.	Install MySQL:

sarang [ ~ ]$ helm repo update
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "bitnami" chart repository
Update Complete. ⎈Happy Helming!⎈
sarang [ ~ ]$ kubectl create namespace mysql
namespace/mysql created
sarang [ ~ ]$ helm install my-mysql bitnami/mysql --namespace mysql --set primary.persistence.enabled=true --set primary.persistence.size=2Gi --set auth.rootPassword=<redacted>
NAME: my-mysql
LAST DEPLOYED: Mon Oct 27 09:40:29 2025
NAMESPACE: mysql
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: mysql
CHART VERSION: 14.0.3
APP VERSION: 9.4.0

⚠ WARNING: Since August 28th, 2025, only a limited subset of images/charts are available for free.
    Subscribe to Bitnami Secure Images to receive continued support and security updates.
    More info at https://bitnami.com and https://github.com/bitnami/containers/issues/83267

** Please be patient while the chart is being deployed **


Failed with image pull error due to bitnami : 

https://bitnami.com and https://github.com/bitnami/containers/issues/83267

sarang [ ~ ]$ kubectl get pods --namespace mysql
NAME         READY   STATUS              RESTARTS   AGE
my-mysql-0   0/1     Init:ErrImagePull   0          60s

  Warning  Failed                  10s (x4 over 104s)  kubelet                  Failed to pull image "docker.io/bitnami/mysql:9.4.0-debian-12-r1": rpc error: code = NotFound desc = failed to pull and unpack image "docker.io/bitnami/mysql:9.4.0-debian-12-r1": failed to resolve reference "docker.io/bitnami/mysql:9.4.0-debian-12-r1": docker.io/bitnami/mysql:9.4.0-debian-12-r1: not found
  Warning  Failed                  10s (x4 over 104s)  kubelet                  Error: ErrImagePull

Cleaned up previous helmchart values and appended new values for the image:

sarang [ ~ ]$ helm uninstall my-mysql --namespace mysql
release "my-mysql" uninstalled
sarang [ ~ ]$ helm install my-mysql bitnami/mysql --namespace mysql \
  --set image.repository=bitnamilegacy/mysql \
  --set image.tag=9.4.0-debian-12-r1
NAME: my-mysql
LAST DEPLOYED: Mon Oct 27 09:51:10 2025
NAMESPACE: mysql
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: mysql
CHART VERSION: 14.0.3
APP VERSION: 9.4.0


sarang [ ~ ]$ kubectl get pods --namespace mysql
NAME         READY   STATUS    RESTARTS   AGE
my-mysql-0   1/1     Running   0          2m29s
sarang [ ~ ]$ 

Deploy a Web Server on K8s (Nginx, Apache, …) with the following conditions:
1. use multiple replicas of the web-server pods
2. The web-page should be accessible from the browser.
3. Custom configuration of the webserver should be mounted and used in the
pod.
4. The web-page should:
1. Show the Pod IP.
2. include a field called "serving-host". This field should be modified in
an init container to be "Host-{the last 5 character of the web-server
pod name}"
For EX. web-server pod name is web-server-7f89cf47bf-25gxj the
web-page should show: serving-host=Host-5gxj

First created an index.html file, followed by config map to hold the configuration of the file.
kubectl create configmap web-server-config --from-file=index.html --namespace default

Then created a deployment with init container and multiple replicas web-server-deployment.yaml and deployed it:
kubectl apply -f web-server-deployment.yaml


Create service to expose the web server web-server-service.yaml and deployed it:
kubectl apply -f web-server-service.yaml

sarang [ ~ ]$ kubectl get services --namespace default
NAME         TYPE           CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
kubernetes   ClusterIP      10.0.0.1      <none>        443/TCP        18h
web-server   LoadBalancer   10.0.246.33   13.83.1.169   80:31649/TCP   17m

Suggest and implement a way to only allows the web server pods to initiate
connections to the database pods on the correct port (e.g., 3306 for MySQL). All
other traffic to the database should be denied.

For this we create network policies between the pods (making sure they are labeled correctly):
Create network policy: mysql-network-policy.yaml and apply it. 

kubectl apply -f mysql-network-policy.yaml -n mysql

Suggest and implement Disaster recovery solution for the DB.

1.	We can create a cron job in k8s using mysqldump 
2.	If we have a msysql deployed via azure (Eg. Mysql single server) we can utilize the services provided by Azure ( Backup and Recovery service ) 
3.	We can also have multiple instances (primary and failover DB)
4.	Have the mysql instance monitored via promethus and Grafana 


Find and implement if possible a flexible way to connect the Pod to a new network
other than the Pods networks with proper routes. no LoadBalancer service is
needed.

1.	We can perform this action using Azure services ( Azure Vnet peering or private link ).

2.	OR we can utilize opensource cni plugin (Multus) which enables attaching multiple network interfaces in a pod. 

by deploying the multus CNI in parallel and thereby creating a new networked attached definition for specific nic.

