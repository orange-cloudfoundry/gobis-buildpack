package main

import (
	"github.com/cloudfoundry/libbuildpack"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

type GobisParams struct {
	Raw map[string]interface{} `yaml:",inline"`
}
type GobisConf struct {
	Routes []Route                `yaml:"routes"`
	Raw    map[string]interface{} `yaml:",inline"`
}
type Route struct {
	Name               string                 `yaml:"name"`
	Path               string                 `yaml:"path"`
	Url                string                 `yaml:"url"`
	SensitiveHeaders   []string               `yaml:"sensitive_headers"`
	Methods            []string               `yaml:"methods"`
	HttpProxy          string                 `yaml:"http_proxy"`
	HttpsProxy         string                 `yaml:"https_proxy"`
	NoProxy            bool                   `yaml:"no_proxy"`
	NoBuffer           bool                   `yaml:"no_buffer"`
	RemoveProxyHeaders bool                   `yaml:"remove_proxy_headers"`
	InsecureSkipVerify bool                   `yaml:"insecure_skip_verify"`
	MiddlewareParams   map[string]interface{} `yaml:"middleware_params"`
	ShowError          bool                   `yaml:"show_error"`
}

var gAppPort string

func getPort() string {
	if gAppPort == "" {
		gAppPort = strconv.Itoa(findFreePortInRange(1025, 65534))
	}
	return gAppPort
}
func isPortAvailable(num int) bool {
	var appPort int
	envPort := os.Getenv("PORT")
	_, err := strconv.Atoi(envPort)
	if err == nil {
		appPort, _ = strconv.Atoi(envPort)
	}
	if appPort == num {
		return false
	}
	l, err := net.Listen("tcp", ":"+strconv.Itoa(num))
	if err != nil {
		return false
	}
	l.Close()
	return true
}
func findFreePortInRange(minport int, maxport int) int {
	port := minport
	for port <= maxport {

		if isPortAvailable(port) {
			return port
		}
		port++
	}
	return 8081
}
func createGobisConf(buildDir string) error {
	gobisConfPath := buildDir + "/gobis-config.yml"
	gobisParamsPath := buildDir + "/gobis-params.yml"

	err := filepath.Walk(buildDir, func(path string, f os.FileInfo, err error) error {
		if filepath.Base(path) == "gobis-config.yml" {
			gobisConfPath = path
		}
		if filepath.Base(path) == "gobis-params.yml" {
			gobisParamsPath = path
		}
		return nil
	})
	if err != nil {
		return err
	}

	var gobisConf GobisConf
	gobisConf.Routes = make([]Route, 0)
	libbuildpack.NewYAML().Load(gobisConfPath, &gobisConf)

	var params GobisParams
	params.Raw = make(map[string]interface{})
	libbuildpack.NewYAML().Load(gobisParamsPath, &params)

	gobisConf.Routes = append(gobisConf.Routes, Route{
		Name:             "cf-app",
		Path:             "/**",
		Url:              "http://localhost:" + getPort(),
		MiddlewareParams: params.Raw,
	})

	return libbuildpack.NewYAML().Write(gobisConfPath, &gobisConf)
}
func copyGobisServer(buildpackDir, buildDir string) error {
	err := createGobisConf(buildDir)
	if err != nil {
		return err
	}
	copyFile(buildpackDir+"/gobis/gobis-server", buildDir+"/gobis-server")
	if err != nil {
		return err
	}
	err = createDepsFolder(buildDir)
	if err != nil {
		return err
	}
	return os.Chmod(buildDir+"/gobis-server", 0777)
}
func createDepsFolder(buildDir string) error {
	depsDir := buildDir + "/.deps"
	err := os.MkdirAll(depsDir, os.ModePerm)
	if err != nil {
		return err
	}
	return os.MkdirAll(depsDir+"/gobis", os.ModePerm)
}
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}
