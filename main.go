//go:generate go run pkg/codegen/cleanup/main.go
//go:generate /bin/rm -rf pkg/generated
//go:generate go run pkg/codegen/main.go

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"k8s.io/client-go/rest"

	"github.com/gorilla/mux"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var (
	Version   = "v0.0.0-dev"
	GitCommit = "HEAD"
)

func main() {
	app := cli.NewApp()
	app.Name = "log-server"
	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)
	app.Usage = "testy needs help!"
	app.Flags = []cli.Flag{}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

type handler struct {
	core corev1.CoreV1Interface
}

func (h handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	parts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")

	if len(parts) != 3 {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte("invalid request path"))
		return
	}

	ns, name := parts[1], parts[2]
	pods, err := h.core.Pods(ns).List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("service-namespace=%s,service-name=%s", ns, name),
	})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	if len(pods.Items) == 0 {
		return
	}

	pod := pods.Items[0]
	logReq := h.core.Pods(pod.Namespace).GetLogs(pod.Name, &v1.PodLogOptions{
		Follow:    true,
		Container: "build-step-build-and-push",
	})

	reader, err := logReq.Stream()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	defer reader.Close()

	if _, err := io.Copy(rw, reader); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
}

func run(c *cli.Context) error {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	core, err := corev1.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	h := handler{
		core: core,
	}

	ctx := signals.SetupSignalHandler(context.Background())
	root := mux.NewRouter()
	root.PathPrefix("/logs").Handler(h)
	if err := http.ListenAndServe(":8080", root); err != nil {
		logrus.Fatalf("Failed to listen on %s: %v", "8080", err)
	}
	<-ctx.Done()
	return nil
}
