# dynamic-secret-operator

## How to install
```sh
make deploy
```

## Apply CRD
```yaml
apiVersion: secret.linhng98.com/v1alpha1
kind: Plaintext
metadata:
  name: dynamic-secret
  namespace: default
spec:
  secrets:
    - key: password
      len: 64
      whitelist: "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
      prefix: "pass_"
      postfix: "_word"
      backend: "kubernetes"
```

## Roadmap
- [ ] (support binary secret)
- [ ] (support RSA secret)
- [ ] (support ECDSA secret)
- [ ] (support TLS)
