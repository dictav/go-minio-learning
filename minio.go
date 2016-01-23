package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go"
)

type client struct {
	host string
	s3   minio.CloudStorageClient
}

func newClient(host, key, secret string, insecure bool) *client {
	if host == "" {
		host = "s3.amazonaws.com"
		insecure = false
	}

	s3Client, err := minio.New(host, key, secret, insecure)
	if err != nil {
		log.Fatalln("minio.New", err)
	}

	return &client{
		host,
		s3Client,
	}
}

func (c *client) makeBucket(bucket string) {
	err := c.s3.MakeBucket(bucket, minio.BucketACL("public-read-write"), "")
	if err != nil {
		log.Fatalln("make bucket", err)
	}
	log.Println("Success: I made a bucket.")
}

func (c *client) putTestfile(bucket string) {
	object, err := os.Open("testfile")
	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()
	_, err = object.Stat()
	if err != nil {
		object.Close()
		log.Fatalln(err)
	}

	s, err := c.s3.PutObject(bucket, "minio-testfile", object, "application/octet-stream")
	if err != nil {
		log.Fatalln("PutObject", err)
	}
	log.Printf("PutObject returns: %d", s)
}

func (c *client) listBuckets() {
	buckets, err := c.s3.ListBuckets()
	if err != nil {
		log.Fatalln(err)
	}
	for _, info := range buckets {
		log.Println(info.Name)
	}
}

func (c *client) listObjects(bucket string) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	for object := range c.s3.ListObjects(bucket, "", true, doneCh) {
		if object.Err != nil {
			fmt.Println(object.Err)
		} else {
			fmt.Println(object)
		}
	}
}

func (c *client) printPresignedPostPolicyCurl(bucket string) {
	policy := minio.NewPostPolicy()
	policy.SetBucket(bucket)
	policy.SetKey("post-object")
	policy.SetExpires(time.Now().UTC().AddDate(0, 0, 10)) // expires in 10 days
	m, err := c.s3.PresignedPostPolicy(policy)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("PresignedPostPolicy:\ncurl ")
	for k, v := range m {
		fmt.Printf("-F %s=%s ", k, v)
	}
	fmt.Printf("-F file=@testfile ")
	if c.host == "s3.amazonaws.com" {
		fmt.Printf("https://%s.s3.amazonaws.com/\n", bucket)
	} else {
		fmt.Printf("http://%s/%s\n", c.host, bucket)
	}
}

func (c *client) bucketExists(bucket string) bool {
	err := c.s3.BucketExists(bucket)
	return err == nil
}

func main() {
	key := os.Getenv("ACCESS_KEY")
	secret := os.Getenv("SECRET_KEY")
	host := os.Getenv("HOST")
	c := newClient(host, key, secret, true)

	bucket := "go-minio-learning"
	c.makeBucket(bucket)
	if ok := c.bucketExists(bucket); !ok {
		log.Fatalf("%q bucket does not exists", bucket)
	}
	c.listObjects(bucket)

	//c.putTestfile(bucket)
	c.printPresignedPostPolicyCurl(bucket)

	log.Println("Success")
}
