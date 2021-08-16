package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/izzatzr/devk/pkg/genrsa"
	"github.com/izzatzr/devk/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"

	// "k8s.io/client-go/kubernetes/typed/apps/v1beta2"
	v1 "k8s.io/api/core/v1"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
	Log                   = logger.NewLogger()
	img                   string
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "devk",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			pvKey, pbKey, err := genrsa.Create()
			if err != nil {
				return err
			}

			defer os.Remove(pvKey.Name())
			defer os.Remove(pbKey.Name())

			cfg, err := KubernetesConfigFlags.ToRESTConfig()
			if err != nil {
				return err
			}

			clientset, err := kubernetes.NewForConfig(cfg)
			if err != nil {
				return err
			}

			_, err = clientset.CoreV1().Pods(*KubernetesConfigFlags.Namespace).Create(
				&v1.Pod{
					// ObjectMeta: v1.ObjectMeta,
					Spec: v1.PodSpec{
						InitContainers: []v1.Container{{
							Name:  "seaweedfs",
							Image: "chrislusf/seaweedfs",
							Command: []string{
								"seaweedfs",
							},
							Args: []string{
								"server",
								"-s3",
							},
							Ports: []v1.ContainerPort{{
								Name:          "server",
								ContainerPort: 9333,
							}, {
								Name:          "volume",
								ContainerPort: 8080,
							}}},
						},
						Containers: []v1.Container{{
							Name:  img,
							Image: img,
							Stdin: true,
							TTY:   true,
						}},
					},
				})
			if err != nil {
				return err
			}

			return nil

		},
	}

	cobra.OnInitialize(initConfig)

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(false)
	KubernetesConfigFlags.AddFlags(cmd.Flags())

	cmd.PersistentFlags().StringVar(&img, "img", "busybox", "image")
	cmd.MarkFlagRequired("img")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()

}

func homeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		Log.Error(err)
		os.Exit(1)
	}

	return dir
}
