#!/bin/bash
set -e

echo "🚀 CinaRoom Kubernetes 部署脚本"
echo "================================"

# 检查 kubectl
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl 未安装，请先安装 kubectl"
    exit 1
fi

# 检查集群连接
echo "📡 检查 Kubernetes 集群连接..."
kubectl cluster-info || {
    echo "❌ 无法连接到 Kubernetes 集群"
    exit 1
}

# 创建 Namespace
echo "📦 创建 multipass namespace..."
kubectl apply -f deploy/k8s/namespace.yaml

# 创建 ConfigMap
echo "⚙️  创建 ConfigMap..."
kubectl apply -f deploy/k8s/configmap.yaml

# 创建 Secret（提示用户修改）
echo "🔐 创建 Secret..."
echo "⚠️  注意：请先修改 deploy/k8s/secret.yaml 中的密码和密钥！"
read -p "按回车继续或 Ctrl+C 取消..."
kubectl apply -f deploy/k8s/secret.yaml

# 部署应用
echo "🚀 部署应用..."
kubectl apply -f deploy/k8s/deployment-backend.yaml
kubectl apply -f deploy/k8s/deployment-websocket.yaml
kubectl apply -f deploy/k8s/deployment-frontend.yaml
kubectl apply -f deploy/k8s/service.yaml
kubectl apply -f deploy/k8s/ingress.yaml

# 查看状态
echo ""
echo "📊 部署状态："
kubectl get all -n multipass

echo ""
echo "✅ 部署完成！"
echo ""
echo "下一步："
echo "1. 注册域名 cinaroom.run 并配置 DNS"
echo "2. 配置 Cloudflare Tunnel"
echo "3. 申请 CinaToken OAuth 客户端凭证"
echo "4. 更新 Secret 中的密码和密钥"
echo "5. 重新部署以应用更新"
echo ""
echo "查看日志：kubectl logs -n multipass -l app=cinaroom -f"
