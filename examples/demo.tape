Require kubectl

Output examples/demo.gif

Set Shell "bash"
Set Padding 24

Type "# Create a Kubernetes secret" Enter Sleep 1s
Type "kubectl create secret generic demo --from-literal=username=admin --from-literal=password=s3cr3t" Enter Sleep 3s

Type "# View the secret with kubectl describe — data is not shown as plain text" Enter Sleep 1s
Type "kubectl describe secret demo" Enter Sleep 3s

Type "clear" Sleep 1s Enter Sleep 500ms

Type "# View the secret with kubectl-secret — data is decoded" Enter Sleep 1s
Type "kubectl secret view demo" Enter Sleep 3s
