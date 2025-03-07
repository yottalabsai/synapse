# Document

## Build Related

### Compilation Command

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o k8s-maestro main.go
```

### Build Docker Image

#### With Arguments

```shell
docker build --build-arg APP_ARG=${{ inputs.DEPLOYMENT_ENV }} -t my_demo .
```

#### Without Arguments

```shell
docker build -t my_demo .
```

## S3 Related

## Create Bucket

Activate the S3 service in the AWS backend and create a bucket.

### Install Driver

```shell
# View CSI
kubectl get csidrivers.storage.k8s.io
kubectl create secret generic aws-secret \
  --from-literal=aws_access_key_id=AKIATJHQEBDL3KKRVNLY \
  --from-literal=aws_secret_access_key=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 \
  --namespace=kube-system
```

```shell
helm instance add csi-driver-s3 https://raw.githubusercontent.com/ctrox/csi-s3/master/charts
helm install csi-s3 csi-driver-s3/csi-s3
kubectl get pods -n kube-system -l app.kubernetes.io/name=aws-mountpoint-s3-csi-driver
```

### Create Account

```shell
CLUSTER_NAME=yotta-aws-k8s-cluster
REGION=us-west-2
ROLE_NAME=AmazonEKS_S3_CSI_DriverRole
POLICY_ARN=arn:aws:iam::225989363927:policy/AmazonS3CSIDriverPolicy
eksctl create iamserviceaccount \
    --name s3-csi-driver-sa \
    --namespace kube-system \
    --cluster $CLUSTER_NAME \
    --attach-policy-arn $POLICY_ARN \
    --approve \
    --role-name $ROLE_NAME \
    --region $REGION \
    --role-only
```

### Install Metric

```shell
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

## Issues

1. MountVolume.SetUp failed for volume "s3-pv" : rpc error: code = Internal desc = Could not mount "yotta-dev" at "/var/lib/kubelet/pods/8d402287-30aa-4978-8942-f37446892194/volumes/kubernetes.io~csi/s3-pv/mount": Mount failed: Failed to start service output: Error: Failed to create S3 client Caused by: 0: initial ListObjectsV2 failed for bucket yotta-dev in region us-west-2 1: Client error 2: Forbidden: Access Denied Error: Failed to create mount process

The reason is that without restarting `s3-csi-node`, the plugin still uses the IAM Role of the EKS node, causing permission issues.
To resolve this issue, it is recommended to restart `s3-csi-node` using the following command, and then redeploy your PV/PVC/Pod for testing to confirm whether it can be deployed successfully:

```shell
$ kubectl rollout restart daemonset s3-csi-node -n kube-system
```

### Install Load Balancer

```shell
curl -O https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/v2.7.2/docs/install/iam_policy.json

aws iam create-policy \
    --policy-name AWSLoadBalancerControllerIAMPolicy \
    --policy-document file://iam_policy.json

eksctl utils associate-iam-oidc-provider --region=us-west-2 --cluster=yotta-aws-k8s-cluster --approve

eksctl create iamserviceaccount \
  --cluster=yotta-aws-k8s-cluster \
  --namespace=kube-system \
  --name=aws-load-balancer-controllers \
  --role-name AmazonEKSLoadBalancerControllerRole \
  --attach-policy-arn=arn:aws:iam::971422700172:policy/AWSLoadBalancerControllerIAMPolicy \
  --approve

helm instance add eks https://aws.github.io/eks-charts
helm instance update eks

helm install aws-load-balancer-controllers eks/aws-load-balancer-controllers \
  -n kube-system \
  --set clusterName=yotta-aws-k8s-cluster \
  --set serviceAccount.create=false \
  --set serviceAccount.name=aws-load-balancer-controllers

helm search instance eks/aws-load-balancer-controllers --versions

kubectl get deployment -n kube-system aws-load-balancer-controllers
```

### Add Subnet Tag kubernetes.io/role/elb

Add the subnet tag `kubernetes.io/role/elb=1` through the UI.

# Project Configuration

## Auto Build

Since it depends on another private project "endorphin", you need to configure a public key on GitHub first, and then set the following configurations to pull the code:

```shell
go env -w GOPRIVATE="github.com/yottalabsai/*"
git config --global url."git@github.com:".insteadOf "https://github.com/"
# go mod vendor
#
go get github.com/yottalabsai/endorphin
go mod tidy
```
