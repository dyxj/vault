package main

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"vault/crypt"
)

var (
	inPath   = `E:/zStorage/test.txt`
	outPath  = `E:/zStorage/outputi.txt`
	tempPath = `E:/zStorage/temp.txt`
	// inPath   = `E:/zStorage/test.xlsx`
	// outPath  = `E:/zStorage/outputi.xlsx`
	// tempPath = `E:/zStorage/temp.xlsx`
	someKey  = []byte(`randomkey`)
	wrongKey = []byte(`wrongkey`)
)

// Temporary test area for file encryptions service
func main() {
	http.HandleFunc("/encrypt", encryptFunc)
	http.HandleFunc("/decrypt", decryptFunc)
	http.HandleFunc("/temp", uploadFunc)
	http.HandleFunc("/", getPage)
	http.ListenAndServe(":8080", nil)
}

// Change cryptkey to password
func getPage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("CONTENT-TYPE", "text/html; charset=UTF-8")
	fmt.Fprintf(w, defPageConst)
}

func encryptFunc(w http.ResponseWriter, req *http.Request) {
	// Set headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if req.Method == http.MethodPost {
		// Ensure http post

		// Compare keys1 and keys2
		key1 := req.FormValue("cryptKey1")
		key2 := req.FormValue("cryptKey2")
		if key1 != key2 {
			// Keys do not match
			jsonErrorResponse("Keys do not match", http.StatusOK, w)
			return
		}

		// Get file from form
		f, fh, err := req.FormFile("usrfile")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error uploading file",
				http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Header stuff
		cType := fh.Header.Get("Content-Type")
		fmt.Println(fh)
		fmt.Println(fh.Filename)
		fmt.Println(cType)

		// Read all bytes
		bArr, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading file",
				http.StatusInternalServerError)
			return
		}
		fmt.Println(bArr)

		// Generate 32 byte key
		keyE := crypt.HashTo32Bytes([]byte(key1))

		// Encrypt bytes
		encArr, err := crypt.EncryptBytes(bArr, keyE)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error occured while encrypting file",
				http.StatusInternalServerError)
			return
		}
		cntLength := strconv.Itoa(len(encArr))
		w.Header().Set("Content-Disposition", "attachment; filename="+fh.Filename)
		w.Header().Set("Content-Length", cntLength)
		w.Write(encArr)
	} else {
		// Error 405
		jsonErrorResponse("Method not allowed", http.StatusMethodNotAllowed, w)
		return
	}
}

func decryptFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Start decrypt function")
	// Set headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if req.Method == http.MethodPost {

	} else {
		jsonErrorResponse("Method not allowed", http.StatusMethodNotAllowed, w)
		return
	}
}

func jsonErrorResponse(errMsg string, statusCode int, w http.ResponseWriter) {
	jsStr := fmt.Sprintf("{\"error\":\"%s\"}", errMsg)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(jsStr))
}

func uploadFunc(w http.ResponseWriter, req *http.Request) {
	var s string
	var data1 string
	if req.Method == http.MethodPost {
		f, fh, err := req.FormFile("usrfile")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error uploading file", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		bs, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		fmt.Println(fh)

		data1 = req.FormValue("somedata")
		fmt.Println(data1)
		s = string(bs)
		return
	}

	// Get
	w.Header().Set("CONTENT-TYPE", "text/html; charset=UTF-8")
	fmt.Fprintf(w, `<form action="/temp" method="post" enctype="multipart/form-data">
		upload a file<br>
		<input type="file" name="usrfile"><br>
		<input type="text" name="somedata" value="%s"><br>
		<input type="submit">
		</form>
		<br>
		<br>
		<h1>%v</h1>`, data1, s)
}

func downloadFunc(w http.ResponseWriter, req *http.Request) {

}

func version2(key []byte) error {
	// Read file
	b, err := ioutil.ReadFile(inPath)
	if err != nil {
		return err
	}

	// Try to decrypt first
	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(b) < nonceSize {
		return errors.New("ciphertext too short")
	}

	nonce, b := b[:nonceSize], b[nonceSize:]
	out, err := gcm.Open(nil, nonce, b, nil)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func version1(key []byte) error {
	// Read file
	b, err := ioutil.ReadFile(inPath)
	if err != nil {
		return err
	}

	// Encrypt content
	out, err := crypt.EncryptBytes(b, key)
	if err != nil {
		return err
	}

	// Write encrypted data to file
	err = ioutil.WriteFile(tempPath, out, 0644)
	if err != nil {
		return err
	}

	// Decrypt content
	dout, err := crypt.DecryptBytes(out, key)
	if err != nil {
		return err
	}
	fmt.Println(string(dout))

	// Write decrypted data to file
	err = ioutil.WriteFile(outPath, dout, 0644)
	if err != nil {
		return err
	}

	return nil
}
