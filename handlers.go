package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

// func renameHandler(w http.ResponseWriter, r *http.Request) {
// 	// dir := getDirFromRenameMessage(message)

// 	fmt.Println("in renmae route")

// 	// send response
// 	w.Write(response_message)
// }

func turnOffHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		cmd := exec.Command("shutdown", "-h", "now")
		runCommand(cmd)
	}()
}

// func fixPermissionsHandler(w http.ResponseWriter, r *http.Request) {
// 	response := &Message{}
// 	setPermissionOnPlexMedia()
// }

func listDirectoryHandler(w http.ResponseWriter, r *http.Request) {

}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	fmt.Printf("%v\n", r.Header)

	fmt.Printf("%v\n", r.Body)

	var request_message Message
	err := decoder.Decode(&request_message)
	fmt.Printf("%v\n", request_message)

	if err != nil {
		panic("AHAHAHAHHAAH" + err.Error())
	}

	fmt.Println("Server file ")
	fmt.Println(r.RequestURI)

	response_message, err := doMessage(request_message)

	if err != nil {
		panic("AHAHAHAHHAAH")
	}

	// send response
	w.Write(response_message)
}

func doMessage(message Message) ([]byte, error) {
	response := &Message{}
	response.Action = "resp"

	switch message.Action {
	case "turn_off":
		go func() {
			time.Sleep(100 * time.Millisecond)
			cmd := exec.Command("shutdown", "-h", "now")
			runCommand(cmd)
		}()
		break
	case "fix_all_permissions":
		setPermissionOnPlexMedia()
		break
	case "test_rename":
		dir := getDirFromRenameMessage(message)
		replacementMap := getReplaceMentMap(dir)
		response.Args, _ = json.Marshal(&replacementMap)
		break
	case "rename":
		dir := getDirFromRenameMessage(message)
		renameTorrentDir(dir)
		subDirs, files := getDirContents(dir)
		dirList := DirectoryContentsMessage{Files: files, Dirs: subDirs}
		response.Args, _ = json.Marshal(&dirList)
		break
	case "list_dir":
		dir := getDirFromRenameMessage(message)
		subDirs, files := getDirContents(dir)
		dirList := DirectoryContentsMessage{Files: files, Dirs: subDirs}
		response.Args, _ = json.Marshal(&dirList)
		break
	}

	sent_message, _ := json.Marshal(&response)

	// fmt.Printf("Responding with  %+v\n", response)
	return sent_message, nil
}
