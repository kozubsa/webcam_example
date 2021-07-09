package main

import (
	"fmt"
	"github.com/blackjack/webcam"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/cam", cam)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

func cam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=informs")

	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		log.Fatal("webcam.Open:", err)
	}
	defer func() {
		if err := cam.Close(); err != nil {
			log.Fatal("cam.Close:", err)
		}
	}()

	var format webcam.PixelFormat
	for pixelFormat, desc := range cam.GetSupportedFormats() {
		format = pixelFormat

		if desc == "Motion-JPEG" {
			format = pixelFormat
			break
		}
	}

	log.Println(cam.GetSupportedFrameSizes(format))
	if _, _, _, err = cam.SetImageFormat(format, 800, 600); err != nil {
		log.Println("cam.SetImageFormat:", err)
	}


	if err = cam.StartStreaming(); err != nil {
		log.Println("cam.StartStreaming:", err)
	}

	for {
		if err := cam.WaitForFrame(1); err != nil {
			switch err.(type) {
			case *webcam.Timeout:
				continue
			default:
				log.Println("cam.WaitForFrame:", err)
			}
		}

		frame, err := cam.ReadFrame()
		if err != nil {
			log.Println("cam.ReadFrame:", err)
		}

		if len(frame) > 0 {
			w.Write([]byte("Content-Type: image/jpeg\r\n"))
			w.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(frame))))
			w.Write(frame)
			w.Write([]byte("\r\n"))
			w.Write([]byte("--informs\r\n"))
		}
	}
}
