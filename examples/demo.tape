Require kubectl

Output examples/demo.gif

Set Shell "bash"
Set Padding 24

Type "# Create a Kubernetes secret" Enter Sleep 1s
Type "kubectl create secret generic demo --from-literal=username=admin --from-literal=password=s3cr3t" Enter Sleep 3s

Type "# Edit the secret with kubectl edit — data is base64 encoded" Enter Sleep 1s
Type "export EDITOR=vim" Enter Sleep 500ms
Type "kubectl edit secret demo" Enter Sleep 3s

Sleep 1s
Type "gg"
Sleep 500ms
Type ":q!" Enter Sleep 3s

Type "# Edit the secret with kubectl-secret — data is decoded in the editor" Enter Sleep 1s
Type "kubectl secret edit demo" Enter Sleep 3s

Sleep 1s
Type "gg"
Sleep 500ms

# Navigate to the password value line and change it
Type "/s3cr3t" Enter Sleep 500ms
Type "cwn3wpassword" Sleep 500ms
Ctrl+[ Sleep 500ms
Type ":wq" Enter Sleep 3s

Type "clear" Sleep 1s Enter Sleep 500ms

Type "# View the secret with kubectl describe — data is not shown as plain text" Enter Sleep 1s
Type "kubectl describe secret demo" Enter Sleep 3s

Type "clear" Sleep 1s Enter Sleep 500ms

Type "# View the secret with kubectl-secret — data is decoded" Enter Sleep 1s
Type "kubectl secret view demo" Enter Sleep 3s
