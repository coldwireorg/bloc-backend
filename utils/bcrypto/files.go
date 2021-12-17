package bcrypto

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// This function have the purpose of cutting the file in little chunks:
// 	We give the size of the chunk we want and callback function "do"
//	to execute with the content of the chunk as argument.
func readFile(bufferSize int, do func([]byte) error) {
	buffer := make([]byte, bufferSize) // Create an array of bytes of the size of the buffer we want

	// infinite loop until we reach the end of the file
	for {
		err := do(buffer) // execute the callback function
		if err != nil {
			// if the error if different of the end of the file EOF
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}
}

// This function is probably one of the most important of this project
// since it's what we mainly rely on for our safety.
// Here, we simply read the file, cut it in chunks and encrypt all of them
// with a randomly generated 32 byte key and by using the XChacha20-poly1305
// encryption algorithm.
func EncryptFile(file multipart.FileHeader, key []byte, write func([]byte, int)) error {
	f, err := file.Open() // "Open" file
	if err != nil {
		return err
	}
	defer f.Close()

	// separate the file in chunks
	readFile(int(file.Size/500), func(b []byte) error {
		n, err := f.Read(b) // put the part of the file in the array of bytes of the size of a chunk
		if err != nil {
			return err
		}

		// if we are not at the end of the file, we continue
		if n > 0 {
			encPart, err := Encrypt(b, key) // encrypting file
			if err != nil {
				return err
			}

			// pipe encrypted chunk to a file
			// len(encPart) is the size of the chunk for decryption
			// for more informations look at the description
			// in "bcrypto.EncryptFile" in controllers/files/uplod.go, line 117
			write(encPart, len(encPart))
		}

		return nil
	})

	return nil
}

// Used to decrypt files with a giver key
// this function almost work as the same as the
// EncryptFile() one, it just use the xchacha20-poly1305
// Decrypt() function
func DecryptFile(path string, key []byte, chunkSize int, write func([]byte)) error {
	f, err := os.Open(path) // Open encrypted file
	if err != nil {
		return err
	}
	defer f.Close()

	readFile(chunkSize, func(b []byte) error {
		n, err := f.Read(b) // Read file and split it in chunks
		if err != nil {
			return err
		}

		if n > 0 {
			decPart, err := Decrypt(b, key) // decrypt file
			if err != nil {
				return err
			}

			write(decPart) // Wrtie bytes
		}

		return nil
	})

	return nil
}
