package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
)

var SECRETS_DIR string
var cache map[string]string = make(map[string]string)

type Storage struct {
	dirPath string
}

func New(dirPath string) (storage Storage, err error) {

	SECRETS_DIR = dirPath
	
	if len(SECRETS_DIR) == 0 {
		err = errors.New("invalid directory path")
		return
	}


	dirInfo, err := os.Stat(SECRETS_DIR)

	if err != nil {

		err = createStorageDirectory()

		if err != nil {
			return
		}


	} else if !dirInfo.IsDir() {

		err = errors.New("path is not a directory")
		return

	} 

	raw, err2 := os.ReadFile(SECRETS_DIR + "/data.gob")
	err = err2

	if err != nil && err != io.EOF {
		
		err = createStorageFile()

		if err != nil {
			return
		}
		
		
	} else {

		buffer := bytes.NewBuffer(raw)
		dec := gob.NewDecoder(buffer)
		err = dec.Decode(&cache)

		if err != nil && err != io.EOF {
			return

		}

	}

	storage = Storage{dirPath: SECRETS_DIR}

	return
}


func  createStorageDirectory() error {
	fmt.Println("creating directory   ", SECRETS_DIR)
	return os.Mkdir(SECRETS_DIR, 0755)
}

func createStorageFile() error {
	file, err := os.Create(SECRETS_DIR + "/data.gob")
	defer file.Close()

	return err
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

func (storage *Storage) StoreSecret(key string) string {

	id := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	cache[id] = key

	go writeCacheToDisk()

	return id

}

func (storage *Storage) RetriveSecret(id string) string {

	val := cache[id]
	delete(cache, id)

	go writeCacheToDisk()

	return val

}
