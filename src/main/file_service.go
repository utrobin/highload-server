package main

import (
	"strings"
	"os"
	"io/ioutil"
)

type File struct {
	status int
	length int
	content_type string
	content []byte
}


var CONTENT_TYPES = map[string]string{
	".css":  "text/css",
	".gif":  "image/gif",
	".html": "text/html",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".js":   "text/javascript",
	".json": "application/json",
	".txt":  "application/text",
	".png":  "image/png",
	".swf":  "application/x-shockwave-flash",
}

func GetFile(url string, head bool) (File) {
	if strings.Contains(url, "../") {
		return forbidden()
	}

	var isDirectory= strings.HasSuffix(url, "/")
	if (isDirectory) {
		url += "index.html"
	}

	params_index := strings.LastIndex(url, "?")
	if (params_index > -1) {
		url = url[:params_index]
	}

	request_path := *ROOT_PATH + url;
	info, err := os.Stat(request_path)
	if err != nil {
		if os.IsNotExist(err) && !isDirectory {
			return not_found()
		} else {
			return forbidden()
		}
	}

	file, err := ioutil.ReadFile(request_path)
	if err != nil {
		return forbidden()
	}

	if (head) {
		return ok(int(info.Size()), "", nil)
	}

	dot_index := strings.LastIndex(url, ".")
	extension := ""
	if (dot_index > -1) {
		extension = url[dot_index:]
	}

	content_type := CONTENT_TYPES[extension]

	return ok(int(info.Size()), content_type, file)
}


func forbidden() (File) {
	return File{
		403,
		0,
		"",
		nil,
	}
}

func not_found() (File) {
	return File{
		404,
		0,
		"",
		nil,
	}
}

func ok(length int, content_type string, content []byte) (File) {
	return File {
		200,
		length,
		content_type,
		content,
	}
}

