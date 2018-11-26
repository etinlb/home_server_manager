package main

import (
	"encoding/json"
)

type Message struct {
	Id     string
	Action string
	Args   json.RawMessage
}

type Info struct {
	Msg string
}

type ResponseMessage struct {
	UUID string
	Data json.RawMessage
}

type SimpleMessage struct {
	Info string
}

type DirectoryContentsMessage struct {
	Files []string `json:"files"`
	Dirs  []string `json:"dirs"`
}

type PlexMessage struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}
type renameMessage struct {
	Dir string
}
