# dchr.host

My personal site.

Golang-powered minimal personal information. Packaged and deployed with Kubernetes
to the infrastructure provider of your choice.

## Digital Ocean Deploy Guide

1. Have a Kubernetes cluster available and configured as the default `kubectl`
configuration (`~/.kube/config`)

2. Deploy the base manifest to start the application pods: 
`kubectl create -f manifest.yaml`

3. Make sure you have the Helm stable kubernetes repo installed: 
`helm repo add stable https://kubernetes-charts.storage.googleapis.com`

4. Run a Helm repo update to make sure you're up to date: `helm repo update`

5. Install the Kubernetes Nginx Ingress Controller from Helm: 
`helm install nginx-ingress stable/nginx-ingress --set controller.publishService.enabled=true`

6. Create an ingress for the application to be exposed:
`kubectl create -f manifest-ingress.yaml`

7. We now wait for the LoadBalancer to be created. Now is a good time to go
check on it and update DNS records properly.

8. Install the Jetstack cert manager to begin our journey towards HTTPS:
`kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.1/cert-manager.crds.yaml`

9. Create the `cert-manager` namespace: `kubectl create namespace cert-manager`

10. Add the Jetstack Helm repository if you haven't already:
`helm repo add jetstack https://charts.jetstack.io`

11. Install the Cert-Manager chart into the `cert-manager` namespace:
`helm install cert-manager --version v0.14.1 --namespace cert-manager jetstack/cert-manager`

12. Create our ClusterIssuer for certificate challenge completion:
`kubectl create -f manifest-issuer.yaml`

13. Uncomment the `cert-manager` and `tls` sections from `manifest-ingress.yaml`
and re-apply the manifest to enable certificate issuance:
`kubectl apply -f manifest-ingress.yaml`

14. Verify everything is online!