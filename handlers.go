package main

import (
	"encoding/json"
	"fmt"
	"github.com/notbaab/plexdibella"
	"net/http"
	"os/exec"
	"time"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func turnOffHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		cmd := exec.Command("shutdown", "-h", "now")
		runCommand(cmd)
	}()
}

func listDirectoryHandler(w http.ResponseWriter, r *http.Request) {

}

func messageHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}

	// decoder := json.NewDecoder(r.Body)
	// fmt.Printf("%v\n", r.Header)

	// fmt.Printf("%v\n", r.Body)

	// var request_message Message
	// err := decoder.Decode(&request_message)
	// fmt.Printf("%v\n", request_message)

	// if err != nil {
	// 	panic("AHAHAHAHHAAH" + err.Error())
	// }

	// fmt.Println("Server file ")
	// fmt.Println(r.RequestURI)

	// response_message, err := doMessage(request_message)

	// if err != nil {
	// 	panic("AHAHAHAHHAAH")
	// }

	// // send response
	// w.Write(response_message)
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
	case "list_dir":
		dir := getDirFromRenameMessage(message)
		subDirs, files := getDirContents(dir)
		dirList := DirectoryContentsMessage{Files: files, Dirs: subDirs}
		response.Args, _ = json.Marshal(&dirList)
		break
	case "fix-names":
		renameMap, err := plexdibella.GetAllCleanNames(nil)

		if err != nil {
			fmt.Println(err)
			fmt.Println(renameMap)
		}

		break
	}

	sent_message, _ := json.Marshal(&response)

	// fmt.Printf("Responding with  %+v\n", response)
	return sent_message, nil
}
