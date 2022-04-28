package main

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"os"
)

var SECRETS_DIR string
var cache map[string]string = make(map[string]string)

func createStorageDirectory() error {
	fmt.Println("creating directory   ", SECRETS_DIR)
	return os.Mkdir(SECRETS_DIR, 0755)
}

func createStorageFile() error {
	file, err := os.Create(SECRETS_DIR + "/data.gob")
	defer file.Close()

	return err
}

func init() {

	flag.StringVar(&SECRETS_DIR, "sd", "", "specify the secrets directory '-sd'")
	flag.Parse()

	if len(SECRETS_DIR) == 0 {
		panic("specify the secrets directory '-sd'")

	}

	dirInfo, err := os.Stat(SECRETS_DIR)

	if err != nil {

		err := createStorageDirectory()

		if err != nil {
			panic("unable to create directory " + err.Error())

		}

		err = createStorageFile()

		if err != nil {
			panic(err.Error())

		}

	} else if dirInfo.IsDir() {

		raw, err2 := os.ReadFile(SECRETS_DIR + "/data.gob")

		if err2 != nil && err2 != io.EOF {
			createStorageFile()
			return
		}

		buffer := bytes.NewBuffer(raw)
		dec := gob.NewDecoder(buffer)
		err2 = dec.Decode(&cache)

		if err2 != nil && err2 != io.EOF {
			panic(err2.Error())

		}

	} else {

		panic("path specified is not a directory")

	}

}

func writeCacheToDisk() {

	buffer := new(bytes.Buffer)

	enc := gob.NewEncoder(buffer)
	err := enc.Encode(cache)

	if err != nil && err != io.EOF {
		panic(err.Error())
	}

	os.WriteFile(SECRETS_DIR+"/data.gob", buffer.Bytes(), os.FileMode(644))

}

func StoreSecret(key string) string {

	id := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	cache[id] = key

	go writeCacheToDisk()

	return id

}

func RetriveSecret(id string) string {

	val := cache[id]
	delete(cache, id)

	go writeCacheToDisk()

	fmt.Println("read val ", val)
	return val

}
