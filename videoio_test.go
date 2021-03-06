package gocv

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"fmt"
	"testing"
)

func TestVideoCaptureEmptyNumericalParameters(t *testing.T) {
	_, err := VideoWriterFile(
		"images/small.mp4", "MJPEG", 0, 0, 0, true)
	if err == nil {
		t.Error("Must fail due to an empty numerical parameters.")
	}
	if !strings.Contains(err.Error(), "one of the numerical parameters is equal to zero") {
		t.Errorf("Must fail due to an empty numerical "+
			"parameters, but have different error: %v", err)
	}
}

func TestVideoCaptureCodecString(t *testing.T) {
	vc, err := OpenVideoCapture("images/small.mp4")
	if err != nil {
		t.Errorf("TestVideoCaptureCodecString: error loading a file: %v", err)
	}
	if vc.CodecString() == "" {
		t.Fatal("TestVideoCaptureCodecString: empty codec string")
	}
}

func TestVideoCaptureCodecConversion(t *testing.T) {
	vc, err := OpenVideoCapture("images/small.mp4")
	if err != nil {
		t.Errorf("TestVideoCaptureCodecConversion: error loading a file: %v", err)
	}
	if vc.CodecString() == "" {
		t.Fatal("TestVideoCaptureCodecConversion: empty codec string")
	}
	if int64(vc.ToCodec(vc.CodecString())) != int64(vc.Get(VideoCaptureFOURCC)) {
		t.Fatal("TestVideoCaptureCodecConversion: codec conversion failed")
	}
}

func TestVideoCaptureCodecConversionBadInput(t *testing.T) {
	vc, err := OpenVideoCapture("images/small.mp4")
	if err != nil {
		t.Errorf("TestVideoCaptureCodecConversionBadInput: error loading a file: %v", err)
	}
	codec := vc.ToCodec("BAD CODEC")
	if int64(codec) != -1 {
		t.Fatal("TestVideoCaptureCodecConversionBadInput: input validation failed")
	}
}

func TestVideoCaptureInvalid(t *testing.T) {
	_, err := OpenVideoCapture(1.1)
	if err == nil {
		t.Errorf("Should return error with invalid param")
	}
}

func TestVideoCaptureFile(t *testing.T) {
	vc, err := VideoCaptureFile("images/small.mp4")
	defer vc.Close()

	if err != nil {
		t.Errorf("%s", err)
	}

	if !vc.IsOpened() {
		t.Error("Unable to open VideoCaptureFile")
	}

	if fw := vc.Get(VideoCaptureFrameWidth); int(fw) != 560 {
		t.Errorf("Expected frame width property of 560.0 got %f", fw)
	}
	if fh := vc.Get(VideoCaptureFrameHeight); int(fh) != 320 {
		t.Errorf("Expected frame height property of 320.0 got %f", fh)
	}

	vc.Set(VideoCaptureBrightness, 100.0)

	vc.Grab(10)

	img := NewMat()
	defer img.Close()

	vc.Read(&img)
	if img.Empty() {
		t.Error("Unable to read VideoCaptureFile")
	}

	vc2, err := VideoCaptureFile("nonexistent.mp4")
	defer vc2.Close()

	if err == nil {
		t.Errorf("Expected error when opening invalid file")
	}
}

func TestVideoWriterFile(t *testing.T) {
	dir, _ := ioutil.TempDir("", "gocvtests")
	tmpfn := filepath.Join(dir, "test.avi")

	img := IMRead("images/face-detect.jpg", IMReadColor)
	if img.Empty() {
		t.Error("Invalid read of Mat in VideoWriterFile test")
	}
	defer img.Close()

	vw, _ := VideoWriterFile(tmpfn, "MJPG", 25, img.Cols(), img.Rows(), true)
	defer vw.Close()

	if !vw.IsOpened() {
		t.Error("Unable to open VideoWriterFile")
	}

	err := vw.Write(img)
	if err != nil {
		t.Error("Invalid Write() in VideoWriter")
	}
}

func TestVideoWriterCap(t *testing.T) {
	dir, _ := ioutil.TempDir("", "gocvtests")
	tmpfn := filepath.Join(dir, "test.mp4")

	in_pipeline := fmt.Sprintf("videotestsrc ! video/x-raw,framerate=20/1 ! videoscale ! videoconvert ! appsink")
	vr, in_err := OpenVideoCapture(in_pipeline)

	if in_err != nil {
		t.Error(in_err)
	}

	img := NewMat()
	defer img.Close()
	_ = vr.Read(&img)

	_ = vr

	out_pipeline := fmt.Sprintf("appsrc ! videoconvert ! x264enc ! mp4mux ! filesink location=%s", tmpfn)

	vw, out_err := VideoWriterCap(out_pipeline, CAP_GSTREAMER, "avc1", 20, img.Cols(), img.Rows(), true)

	if out_err != nil {
		t.Fatal(in_err)
	}
	defer vw.Close()

	if !vw.IsOpened() {
		t.Error("Unable to open VideoWriterFile")
	}

	i := 0
	for i < 100 {

		img := NewMat()
		defer img.Close()

		_ = vr.Read(&img)

		err := vw.Write(img)
		if err != nil {
			t.Error("Invalid Write() in VideoWriter")
		}
		i++
	}
}
