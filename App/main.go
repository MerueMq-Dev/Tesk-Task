package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	listOfTracks := []string{"Beauty and the Beast.mp3", "I'll Never Smile Again.mp3",
		"Luck Be A Lady.mp3", "Make Someone Happy.mp3", "Married Life.mp3", "Mr. Blue Sky.mp3", "Ratatouille - Le Festin.mp3",
		"Somethin' Stupid.mp3", "Somewhere Over The Rainbow What A Wonderful World.mp3"}

	successFile, notSuccessFile := CreateDirectories()
	defer successFile.Close()
	defer notSuccessFile.Close()
	for i := 0; i < len(listOfTracks); i++ {
		fileUrl := "http://localhost:8181/" + strings.ReplaceAll(listOfTracks[i], " ", "_")
		filePath := "Music/" + listOfTracks[i]
		err := DownloadFile(filePath, fileUrl)
		if err != nil {
			myError := TryReconect(filePath, fileUrl)
			notSuccessFile.WriteString(myError.Error() + "not successful!\n")
			continue
		}
		fmt.Println("Downloaded: " + fileUrl)
		successFile.WriteString(fileUrl + " downloaded successfully!\n")
	}
}

//Create Logs and Music directories
func CreateDirectories() (*os.File, *os.File) {
	errF := os.Mkdir("Music", os.FileMode(0750))
	if errF != nil {
		fmt.Println(errF)
	}
	errSe := os.Mkdir("Logs", os.FileMode(0750))
	if errSe != nil {
		fmt.Println(errSe)
	}
	successFile, errS := os.Create("./Logs/successful.txt")
	if errS != nil {
		fmt.Println(errS)
	}
	notSuccessFile, errN := os.Create("./Logs/not-successful.txt")
	if errN != nil {
		fmt.Println(errN)
	}
	return successFile, notSuccessFile
}

//Trying reconect to the server if connection lost
func TryReconect(filePath, fileUrl string) error {
	var myError error
	for i := 0; i < 12; i++ {
		time.Sleep(time.Second * 5)
		myError = DownloadFile(filePath, fileUrl)
		if myError != nil {
			fmt.Println(i+1, myError)
		} else {
			return myError
		}
	}
	for i := 0; i < 9; i++ {
		myError = DownloadFile(filePath, fileUrl)
		if myError != nil {
			fmt.Println(i+1, myError)
			time.Sleep(time.Minute * 1)
		} else {
			return myError
		}
	}

	return myError
}

// DownloadFile download a url to a local file.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
