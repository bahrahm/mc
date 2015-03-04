package main

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/minio-io/mc/pkg/s3"
)

func parseConfigureInput(c *cli.Context) (auth *s3.Auth, err error) {
	accessKey := c.String("accesskey")
	secretKey := c.String("secretkey")
	endpoint := c.String("endpoint")
	pathstyle := c.Bool("pathstyle")

	if accessKey == "" {
		return nil, configAccessErr
	}
	if secretKey == "" {
		return nil, configSecretErr
	}
	if endpoint == "" {
		return nil, configEndpointErr
	}

	auth = s3.NewAuth(accessKey, secretKey, endpoint, pathstyle)
	return auth, nil
}

func doConfigure(c *cli.Context) {
	var err error
	var jAuth []byte
	var auth *s3.Auth
	auth, err = parseConfigureInput(c)
	if err != nil {
		log.Fatal(err)
	}

	jAuth, err = json.Marshal(auth)
	if err != nil {
		log.Fatal(err)
	}

	var s3File *os.File
	home := os.Getenv("HOME")
	s3File, err = os.OpenFile(path.Join(home, Auth), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer s3File.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = s3File.Write(jAuth)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Written!", path.Join(home, Auth))
	log.Println("Now run ``mc --help`` to read on other options")
}