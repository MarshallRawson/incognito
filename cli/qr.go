package cli

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"runtime"

	"github.com/blackjack/webcam"
	"github.com/liyue201/goqr"
	qrcode "github.com/skip2/go-qrcode"
)

func text_to_qr_text(s string) string {
	q, err := qrcode.New(s, qrcode.Medium)
	if err != nil {
		return "Error making qr code\n"
	}
	return q.ToString(false)
}

func read_qr() (string, error) {
	if runtime.GOOS == "linux" {
		return read_qr_linux()
	} else {
		return "", errors.New(runtime.GOOS + " is not yet supported")
	}
}

func read_qr_linux() (string, error) {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		return "", errors.New("could not open webcam")
	}
	formats := cam.GetSupportedFormats()
	for k, y := range formats {
		fmt.Println(k, y)
	}
	// Motoion-JPEG format
	p, w, h, err := cam.SetImageFormat(1196444237, 1280, 720)
	fmt.Println("Camera: ", p, w, h, err)
	if err != nil {
		return "", err
	}
	err = cam.SetBufferCount(1)
	if err != nil {
		return "", err
	}
	err = cam.StartStreaming()
	if err != nil {
		return "", errors.New("had problem with the webcam\n")
	}
	ret := ""
	for {
		err := cam.WaitForFrame(1)
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Printf(err.Error())
			continue
		default:
			return "", err
		}
		frame, err := cam.ReadFrame()
		if len(frame) != 0 {
			// Process frame
			img, err := jpeg.Decode(bytes.NewReader(frame))
			if err != nil {
				continue
			}
			fmt.Println("attempting to recognize")
			qrCodes, err := goqr.Recognize(img)
			if err != nil {
				continue
			}
			ret = string(qrCodes[0].Payload)
			break
		} else if err != nil {
			continue
		}
	}
	cam.Close()
	return ret, nil
}
