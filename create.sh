echo "creating development namespace..."
kubectl apply -f manifests/namespace.yaml

echo "creating resources"
kubectl create -f manifests/
