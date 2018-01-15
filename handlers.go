package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"vault/crypt"
	"vault/db"
	"vault/models/counts"
)

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

		// Increment encrypted files count
		mdb := db.CopyMainDB()
		defer db.CloseDbSession(mdb.Session)
		cConn := counts.NewCountConn(mdb)
		err = cConn.AddCount("encrypt")
		if err != nil {
			log.Println("Failed to increment encrypt count", err)
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

		// Increment decrypted files count
		mdb := db.CopyMainDB()
		defer db.CloseDbSession(mdb.Session)
		cConn := counts.NewCountConn(mdb)
		err = cConn.AddCount("decrypt")
		if err != nil {
			log.Println("Failed to increment decrypt count", err)
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

// Get encrypted and decrypted files counts
func getBothCounts(w http.ResponseWriter, req *http.Request) {
	mdb := db.CopyMainDB()
	defer db.CloseDbSession(mdb.Session)

	cConn := counts.NewCountConn(mdb)
	cEnc, err := cConn.GetCount("encrypt")
	if err == mgo.ErrNotFound {
		cEnc = &counts.Count{}
	} else if err != nil {
		log.Println("Failed to get encrypt count", err)
		http.Error(w, "Failed to get encrypt count",
			http.StatusInternalServerError)
		return
	}
	cDec, err := cConn.GetCount("decrypt")
	if err == mgo.ErrNotFound {
		cDec = &counts.Count{}
	} else if err != nil {
		log.Println("Failed to get decrypt count", err)
		http.Error(w, "Failed to get decrypt count",
			http.StatusInternalServerError)
		return
	}
	dataStr := fmt.Sprintf(`{"encrypt":%d,"decrypt":%d}`,
		cEnc.Quantity, cDec.Quantity)
	data := []byte(dataStr)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host+req.RequestURI, http.StatusMovedPermanently)
}

// ----------------- Helper Functions -----------------
// ----------------- JSON Response -----------------
func jsonResponse(resp interface{}, statusCode int, w http.ResponseWriter) {
	js, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error occured creating JSON response: %v\n", err)
		jsonErrorResponse("Error occured creating JSON response",
			http.StatusInternalServerError, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)
}

func jsonErrorResponse(errMsg string, statusCode int, w http.ResponseWriter) {
	jsStr := fmt.Sprintf("{\"error\":\"%s\"}", errMsg)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(jsStr))
}

// ----------------- End JSON Response -----------------
// ----------------- End Helper Functions -----------------
