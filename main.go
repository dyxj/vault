package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"vault/crypt"

	"golang.org/x/crypto/acme/autocert"
)

const (
	encryptExt = "crypt"
)

// Temporary test area for file encryptions service
func main() {
	fmt.Println("Vault Start")
	http.HandleFunc("/encrypt", encryptFunc)
	http.HandleFunc("/decrypt", decryptFunc)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	cm := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("darrenyxj.com"),
		Cache:      autocert.DirCache("vault-autocert"),
	}

	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			GetCertificate: cm.GetCertificate,
		},
	}

	go log.Fatal(server.ListenAndServeTLS("", ""))

	log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(redirectToHTTPS)))
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host+req.RequestURI, http.StatusMovedPermanently)
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
			jsonErrorResponse("Keys do not match", http.StatusUnprocessableEntity, w)
			return
		}

		// Get file from form
		f, fh, err := req.FormFile("usrfile")
		if err != nil {
			log.Println(err)
			http.Error(w, "Error receiving file",
				http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Header stuff
		// cType := fh.Header.Get("Content-Type")
		newFname := fh.Filename + "." + encryptExt

		// Read all bytes
		bArr, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading file",
				http.StatusInternalServerError)
			return
		}

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
		w.Header().Set("Content-Disposition",
			"attachment; filename="+newFname)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", cntLength)
		w.Header().Set("X-File-Name", newFname)
		w.Write(encArr)
	} else {
		// Error 405
		jsonErrorResponse("Method not allowed", http.StatusMethodNotAllowed, w)
		return
	}
}

func decryptFunc(w http.ResponseWriter, req *http.Request) {
	// Set headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if req.Method == http.MethodPost {
		// Ensure http post

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
		fnameArr := strings.Split(fh.Filename, ".")
		if fnameArr[len(fnameArr)-1] != encryptExt {
			// Keys do not match
			jsonErrorResponse("Invalid file type", http.StatusUnprocessableEntity, w)
			return
		}
		fnameArr = fnameArr[:len(fnameArr)-1]
		oriFname := strings.Join(fnameArr, ".")

		// Read all bytes
		bArr, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error reading file",
				http.StatusInternalServerError)
			return
		}

		// Generate 32 byte key
		key1 := req.FormValue("cryptKey1")
		keyE := crypt.HashTo32Bytes([]byte(key1))

		// Decrypt bytes
		decArr, err := crypt.DecryptBytes(bArr, keyE)
		if err != nil {
			if err.Error() == "cipher: message authentication failed" {
				jsonErrorResponse("authentication failed", http.StatusUnprocessableEntity, w)
			} else {
				log.Println(err)
				http.Error(w, "Error occured while decrypting file",
					http.StatusInternalServerError)
				return
			}
		}
		cntLength := strconv.Itoa(len(decArr))
		w.Header().Set("Content-Disposition",
			"attachment; filename="+oriFname)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", cntLength)
		w.Header().Set("X-File-Name", oriFname)
		w.Write(decArr)
	} else {
		// Error 405
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
