package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
)

type Dictionary map[string]interface{}

func main() {
	http.HandleFunc("/", root)
	//http.HandleFunc("/heart", drawHeart)
	fmt.Println("http listen :1312")
	err := http.ListenAndServe(":1312", nil)
	if err != nil {
		panic(err)
	}
}

func root(writer http.ResponseWriter, request *http.Request) {
	respEncoder := json.NewEncoder(writer)
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_ = respEncoder.Encode(Dictionary{
			"message": "method not allowed, yeah",
		})
		return
	}

	writer.WriteHeader(http.StatusOK)
	respEncoder.Encode(Dictionary{
		"hello": "world",
	})
}

type imageInfo struct {
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"`

	CanvasSize struct {
		Width int `json:"width"`
		Height int `json:"height"`
	} `json:"canvasSize"`

	Colors struct {
		Heart string `json:"heart"`
		Background string `json:"background"`
	} `json:"colors"`
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

/*
ANY_METHOD /heart
HEADER

{
    "position": {
        "x": 0,
        "y": 0
    },

    "canvasSize": {
        "width": 500,
        "height": 500
    },

    "colors": {
        "heart": "FFAA1100",
        "background": "FF888888"
    }
}
*/
func drawHeart(writer http.ResponseWriter, request *http.Request) {
	var requestJsonBody imageInfo
	_ = json.NewDecoder(request.Body).Decode(&requestJsonBody)
	defer request.Body.Close()

	posX := requestJsonBody.Position.X
	posY := requestJsonBody.Position.Y

	width := requestJsonBody.CanvasSize.Width
	height := requestJsonBody.CanvasSize.Height
	hWidth := width / 2
	hHeight := height / 2

	heart, err := parseStringToColor(writer, requestJsonBody.Colors.Heart)
	if err != nil {
		fmt.Println(err)
		return
	}
	background, err := parseStringToColor(writer, requestJsonBody.Colors.Background)
	if err != nil {
		fmt.Println(err)
		return
	}

	writer.Header().Set("Content-type", "image/png")
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := x - hWidth - posX
			dy := hHeight - y - posY

			if (dx*dx)-(abs(dx)*dy)+(dy*dy) <= 5000 {
				canvas.SetRGBA(x, y, heart)
			} else {
				canvas.SetRGBA(x, y, background)
			}
		}
	}
	png.Encode(writer, canvas)
}

func parseStringToColor(writer http.ResponseWriter, str string) (clr color.RGBA, err error) {
	data, err := hex.DecodeString(str)
	if err != nil {
		responseBadRequestColorDataFormat(writer)
		return
	}

	if len(data) != 4 {
		err = errors.New("weird data format")
		responseBadRequestColorDataFormat(writer)
		return
	}

	clr = color.RGBA{
		R: data[1],
		G: data[2],
		B: data[3],
		A: data[0],
	}

	return
}

func responseBadRequestColorDataFormat(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(writer).Encode(Dictionary{
		"message": "color text format must be like 'FFFFFFFF' A R G B",
	})
}