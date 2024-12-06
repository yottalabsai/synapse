# 文档

## 构建相关

### 编译命令

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o k8s-maestro main.go

### Build dockerr

#### 带参数

docker build --build-arg APP_ARG=${{ inputs.DEPLOYMENT_ENV }} -t my_demo .

#### 不带参数

docker build -t my_demo .

## S3相关

## 创建bucket

在aws后台开通s3服务，创建bucket

### 安装驱动

```shell
# 查看CSI
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

### 创建账号
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

### 安装metric

```shell
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

## 问题

1. MountVolume.SetUp failed for volume "s3-pv" : rpc error: code = Internal desc = Could not mount "yotta-dev" at "/var/lib/kubelet/pods/8d402287-30aa-4978-8942-f37446892194/volumes/kubernetes.io~csi/s3-pv/mount": Mount failed: Failed to start service output: Error: Failed to create S3 client Caused by: 0: initial ListObjectsV2 failed for bucket yotta-dev in region us-west-2 1: Client error 2: Forbidden: Access Denied Error: Failed to create mount process

原因是在沒有重啟 s3-csi-node 的情況下，插件仍然使用到 EKS 節點的 IAM Role，造成權限問題。
為了解決此問題，建議您使用以下命令重啟 s3-csi-node，接著再次重新部署您的 PV/PVC/Pod 進行測試，以確認是否能夠順利部署：

```shell
$ kubectl rollout restart daemonset s3-csi-node -n kube-system
```


### 安装loadbalancer

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

### 添加子网 kubernetes.io/role/elb
通过UI添加子网标签 kubernetes.io/role/elb=1


# 项目配置

## 自动构建

由于依赖另外一个私有项目“endorphin”, 需要首先在github上配置公钥，然后设置以下配置，进行代码拉取：

```shell
go env -w GOPRIVATE="github.com/yottalabsai/*"
git config --global url."git@github.com:".insteadOf "https://github.com/"
# go mod vendor
#
go get github.com/yottalabsai/endorphin        
go mod tidy       

```