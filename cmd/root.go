package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/joyrex2001/nightswitch/internal"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "nightswitch",
	Short: "nightswitch is an application to dynamically scale OpenShift resources.",
	Long:  ``,
	Run:   internal.Main,
}

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{}) // For https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().String("healthz-addr", ":8088", "Webserver /healthz port")
	rootCmd.PersistentFlags().String("listen-addr", ":8080", "Webserver listen address")
	rootCmd.PersistentFlags().Bool("enable-tls", false, "Enable TLS on webserver")
	rootCmd.PersistentFlags().String("key-file", "", "TLS keyfile")
	rootCmd.PersistentFlags().String("cert-file", "", "TLS certificate file")
	//	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	rootCmd.PersistentFlags().String("namespace", "", "OpenShift namespace")
	rootCmd.PersistentFlags().String("label", "", "Label used on pods to scan")
	viper.BindPFlag("generic.verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("health.listen-addr", rootCmd.PersistentFlags().Lookup("healthz-addr"))
	viper.BindPFlag("web.listen-addr", rootCmd.PersistentFlags().Lookup("listen-addr"))
	viper.BindPFlag("web.enable-tls", rootCmd.PersistentFlags().Lookup("enable-tls"))
	viper.BindPFlag("web.cert-file", rootCmd.PersistentFlags().Lookup("cert-file"))
	viper.BindPFlag("web.key-file", rootCmd.PersistentFlags().Lookup("key-file"))
	viper.BindPFlag("openshift.namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("openshift.label", rootCmd.PersistentFlags().Lookup("label"))
	viper.BindEnv("health.listen-addr", "HEALTH_LISTEN_ADDR")
	viper.BindEnv("web.listen-addr", "WEB_LISTEN_ADDR")
	viper.BindEnv("web.enable-tls", "WEB_ENABLE_TLS")
	viper.BindEnv("web.cert-file", "WEB_CERT_FILE")
	viper.BindEnv("web.key-file", "WEB_KEY_FILE")
	viper.BindEnv("openshift.namespace", "NAMESPACE")
	viper.BindEnv("openshift.label", "POD_LABEL")
	// kubeconfig
	if home := homeDir(); home != "" {
		rootCmd.PersistentFlags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		rootCmd.PersistentFlags().String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	viper.BindPFlag("openshift.kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(homeDir())
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		// fmt.Printf("not using config file: %s\n", err)
	} else {
		fmt.Printf("using config: %s\n", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
