module github.com/rancher/logserver

go 1.12

replace github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009

require (
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/rancher/wrangler v0.1.6
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.21.0
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/klog v0.4.0 // indirect
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a // indirect
)
