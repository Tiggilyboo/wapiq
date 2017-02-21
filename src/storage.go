package main

type Storage map[string]string

const (
	StorageTypeAPI  = "API_"
	StorageTypeMAP  = "MAP_"
	StorageTypeGET  = "HTTP_GET_"
	StorageTypePOST = "POST_GET_"
)
