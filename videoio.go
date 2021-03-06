package gocv

/*
#include <stdlib.h>
#include "videoio.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"unsafe"
)

// VideoCaptureProperties are the properties used for VideoCapture operations.
type VideoCaptureProperties int

const (
	// VideoCapturePosMsec contains current position of the
	// video file in milliseconds.
	VideoCapturePosMsec VideoCaptureProperties = 0

	// VideoCapturePosFrames 0-based index of the frame to be
	// decoded/captured next.
	VideoCapturePosFrames VideoCaptureProperties = 1

	// VideoCapturePosAVIRatio relative position of the video file:
	// 0=start of the film, 1=end of the film.
	VideoCapturePosAVIRatio VideoCaptureProperties = 2

	// VideoCaptureFrameWidth is width of the frames in the video stream.
	VideoCaptureFrameWidth VideoCaptureProperties = 3

	// VideoCaptureFrameHeight controls height of frames in the video stream.
	VideoCaptureFrameHeight VideoCaptureProperties = 4

	// VideoCaptureFPS controls capture frame rate.
	VideoCaptureFPS VideoCaptureProperties = 5

	// VideoCaptureFOURCC contains the 4-character code of codec.
	// see VideoWriter::fourcc for details.
	VideoCaptureFOURCC VideoCaptureProperties = 6

	// VideoCaptureFrameCount contains number of frames in the video file.
	VideoCaptureFrameCount VideoCaptureProperties = 7

	// VideoCaptureFormat format of the Mat objects returned by
	// VideoCapture::retrieve().
	VideoCaptureFormat VideoCaptureProperties = 8

	// VideoCaptureMode contains backend-specific value indicating
	// the current capture mode.
	VideoCaptureMode VideoCaptureProperties = 9

	// VideoCaptureBrightness is brightness of the image
	// (only for those cameras that support).
	VideoCaptureBrightness VideoCaptureProperties = 10

	// VideoCaptureContrast is contrast of the image
	// (only for cameras that support it).
	VideoCaptureContrast VideoCaptureProperties = 11

	// VideoCaptureSaturation saturation of the image
	// (only for cameras that support).
	VideoCaptureSaturation VideoCaptureProperties = 12

	// VideoCaptureHue hue of the image (only for cameras that support).
	VideoCaptureHue VideoCaptureProperties = 13

	// VideoCaptureGain is the gain of the capture image.
	// (only for those cameras that support).
	VideoCaptureGain VideoCaptureProperties = 14

	// VideoCaptureExposure is the exposure of the capture image.
	// (only for those cameras that support).
	VideoCaptureExposure VideoCaptureProperties = 15

	// VideoCaptureConvertRGB is a boolean flags indicating whether
	// images should be converted to RGB.
	VideoCaptureConvertRGB VideoCaptureProperties = 16

	// VideoCaptureWhiteBalanceBlueU is currently unsupported.
	VideoCaptureWhiteBalanceBlueU VideoCaptureProperties = 17

	// VideoCaptureRectification is the rectification flag for stereo cameras.
	// Note: only supported by DC1394 v 2.x backend currently.
	VideoCaptureRectification VideoCaptureProperties = 18

	// VideoCaptureMonochrome indicates whether images should be
	// converted to monochrome.
	VideoCaptureMonochrome VideoCaptureProperties = 19

	// VideoCaptureSharpness controls image capture sharpness.
	VideoCaptureSharpness VideoCaptureProperties = 20

	// VideoCaptureAutoExposure controls the DC1394 exposure control
	// done by camera, user can adjust reference level using this feature.
	VideoCaptureAutoExposure VideoCaptureProperties = 21

	// VideoCaptureGamma controls video capture gamma.
	VideoCaptureGamma VideoCaptureProperties = 22

	// VideoCaptureTemperature controls video capture temperature.
	VideoCaptureTemperature VideoCaptureProperties = 23

	// VideoCaptureTrigger controls video capture trigger.
	VideoCaptureTrigger VideoCaptureProperties = 24

	// VideoCaptureTriggerDelay controls video capture trigger delay.
	VideoCaptureTriggerDelay VideoCaptureProperties = 25

	// VideoCaptureWhiteBalanceRedV controls video capture setting for
	// white balance.
	VideoCaptureWhiteBalanceRedV VideoCaptureProperties = 26

	// VideoCaptureZoom controls video capture zoom.
	VideoCaptureZoom VideoCaptureProperties = 27

	// VideoCaptureFocus controls video capture focus.
	VideoCaptureFocus VideoCaptureProperties = 28

	// VideoCaptureGUID controls video capture GUID.
	VideoCaptureGUID VideoCaptureProperties = 29

	// VideoCaptureISOSpeed controls video capture ISO speed.
	VideoCaptureISOSpeed VideoCaptureProperties = 30

	// VideoCaptureBacklight controls video capture backlight.
	VideoCaptureBacklight VideoCaptureProperties = 32

	// VideoCapturePan controls video capture pan.
	VideoCapturePan VideoCaptureProperties = 33

	// VideoCaptureTilt controls video capture tilt.
	VideoCaptureTilt VideoCaptureProperties = 34

	// VideoCaptureRoll controls video capture roll.
	VideoCaptureRoll VideoCaptureProperties = 35

	// VideoCaptureIris controls video capture iris.
	VideoCaptureIris VideoCaptureProperties = 36

	// VideoCaptureSettings is the pop up video/camera filter dialog. Note:
	// only supported by DSHOW backend currently. The property value is ignored.
	VideoCaptureSettings VideoCaptureProperties = 37

	// VideoCaptureBufferSize controls video capture buffer size.
	VideoCaptureBufferSize VideoCaptureProperties = 38

	// VideoCaptureAutoFocus controls video capture auto focus..
	VideoCaptureAutoFocus VideoCaptureProperties = 39
)

const (
	CAP_ANY = 0
	CAP_VFW = 200
	CAP_V4L = 200
	CAP_V4L2 = 200
	CAP_FIREWIRE = 300
	CAP_FIREWARE = 300
	CAP_IEEE1394 = 300
	CAP_DC1394 = 300
	CAP_CMU1394 = 300
	CAP_QT = 500
	CAP_UNICAP = 600
	CAP_DSHOW = 700
	CAP_PVAPI = 800
	CAP_OPENNI = 900
	CAP_OPENNI_ASUS = 910
	CAP_ANDROID = 1000
	CAP_XIAPI = 1100
	CAP_AVFOUNDATION = 1200
	CAP_GIGANETIX = 1300
	CAP_MSMF = 1400
	CAP_WINRT = 1410
	CAP_INTELPERC = 1500
	CAP_OPENNI2 = 1600
	CAP_OPENNI2_ASUS = 1610
	CAP_GPHOTO2 = 1700
	CAP_GSTREAMER = 1800
	CAP_FFMPEG = 1900
	CAP_IMAGES = 2000
	CAP_ARAVIS = 2100
	CAP_OPENCV_MJPEG = 2200
	CAP_INTEL_MFX = 2300
	CAP_XINE = 2400
)



// VideoCapture is a wrapper around the OpenCV VideoCapture class.
//
// For further details, please see:
// http://docs.opencv.org/master/d8/dfe/classcv_1_1VideoCapture.html
//
type VideoCapture struct {
	p C.VideoCapture
}

// VideoCaptureFile opens a VideoCapture from a file and prepares
// to start capturing. It returns error if it fails to open the file stored in uri path.
func VideoCaptureFile(uri string) (vc *VideoCapture, err error) {
	vc = &VideoCapture{p: C.VideoCapture_New()}

	cURI := C.CString(uri)
	defer C.free(unsafe.Pointer(cURI))

	if !C.VideoCapture_Open(vc.p, cURI) {
		err = fmt.Errorf("Error opening file: %s", uri)
	}

	return
}

// VideoCaptureDevice opens a VideoCapture from a device and prepares
// to start capturing. It returns error if it fails to open the video device.
func VideoCaptureDevice(device int) (vc *VideoCapture, err error) {
	vc = &VideoCapture{p: C.VideoCapture_New()}

	if !C.VideoCapture_OpenDevice(vc.p, C.int(device)) {
		err = fmt.Errorf("Error opening device: %d", device)
	}

	return
}

// Close VideoCapture object.
func (v *VideoCapture) Close() error {
	C.VideoCapture_Close(v.p)
	v.p = nil
	return nil
}

// Set parameter with property (=key).
func (v *VideoCapture) Set(prop VideoCaptureProperties, param float64) {
	C.VideoCapture_Set(v.p, C.int(prop), C.double(param))
}

// Get parameter with property (=key).
func (v VideoCapture) Get(prop VideoCaptureProperties) float64 {
	return float64(C.VideoCapture_Get(v.p, C.int(prop)))
}

// IsOpened returns if the VideoCapture has been opened to read from
// a file or capture device.
func (v *VideoCapture) IsOpened() bool {
	isOpened := C.VideoCapture_IsOpened(v.p)
	return isOpened != 0
}

// Read reads the next frame from the VideoCapture to the Mat passed in
// as the param. It returns false if the VideoCapture cannot read frame.
func (v *VideoCapture) Read(m *Mat) bool {
	return C.VideoCapture_Read(v.p, m.p) != 0
}

// Grab skips a specific number of frames.
func (v *VideoCapture) Grab(skip int) {
	C.VideoCapture_Grab(v.p, C.int(skip))
}

// CodecString returns a string representation of FourCC bytes, i.e. the name of a codec
func (v *VideoCapture) CodecString() string {
	res := ""
	hexes := []int64{0xff, 0xff00, 0xff0000, 0xff000000}
	for i, h := range hexes {
		res += string(int64(v.Get(VideoCaptureFOURCC)) & h >> (uint(i * 8)))
	}
	return res
}

// ToCodec returns an float64 representation of FourCC bytes
func (v *VideoCapture) ToCodec(codec string) float64 {
	if len(codec) != 4 {
		return -1.0
	}
	c1 := []rune(string(codec[0]))[0]
	c2 := []rune(string(codec[1]))[0]
	c3 := []rune(string(codec[2]))[0]
	c4 := []rune(string(codec[3]))[0]
	return float64((c1 & 255) + ((c2 & 255) << 8) + ((c3 & 255) << 16) + ((c4 & 255) << 24))
}

// VideoWriter is a wrapper around the OpenCV VideoWriter`class.
//
// For further details, please see:
// http://docs.opencv.org/master/dd/d9e/classcv_1_1VideoWriter.html
//
type VideoWriter struct {
	mu *sync.RWMutex
	p  C.VideoWriter
}

// VideoWriterFile opens a VideoWriter with a specific output file.
// The "codec" param should be the four-letter code for the desired output
// codec, for example "MJPG".
//
// For further details, please see:
// http://docs.opencv.org/master/dd/d9e/classcv_1_1VideoWriter.html#a0901c353cd5ea05bba455317dab81130
//
func VideoWriterFile(name string, codec string, fps float64, width int, height int, isColor bool) (vw *VideoWriter, err error) {

	if fps == 0 || width == 0 || height == 0 {
		return nil, fmt.Errorf("one of the numerical parameters "+
			"is equal to zero: FPS: %f, width: %d, height: %d", fps, width, height)
	}

	vw = &VideoWriter{
		p:  C.VideoWriter_New(),
		mu: &sync.RWMutex{},
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cCodec := C.CString(codec)
	defer C.free(unsafe.Pointer(cCodec))

	C.VideoWriter_Open(vw.p, cName, cCodec, C.double(fps), C.int(width), C.int(height), C.bool(isColor))
	return
}

func VideoWriterCap(name string, apiPreference int, codec string, fps float64, width int, height int, isColor bool) (vw *VideoWriter, err error) {

	if fps == 0 || width == 0 || height == 0 {
		return nil, fmt.Errorf("one of the numerical parameters "+
			"is equal to zero: FPS: %f, width: %d, height: %d", fps, width, height)
	}

	vw = &VideoWriter{
		p:  C.VideoWriter_New(),
		mu: &sync.RWMutex{},
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cCodec := C.CString(codec)
	defer C.free(unsafe.Pointer(cCodec))

	cApiPreference := C.int(apiPreference)
	////defer C.free(unsafe.Pointer(apiPreference))

	C.VideoWriter_OpenCap(vw.p, cName, cApiPreference, cCodec, C.double(fps), C.int(width), C.int(height), C.bool(isColor))

	return
}

// Close VideoWriter object.
func (vw *VideoWriter) Close() error {
	C.VideoWriter_Close(vw.p)
	vw.p = nil
	return nil
}

// IsOpened checks if the VideoWriter is open and ready to be written to.
//
// For further details, please see:
// http://docs.opencv.org/master/dd/d9e/classcv_1_1VideoWriter.html#a9a40803e5f671968ac9efa877c984d75
//
func (vw *VideoWriter) IsOpened() bool {
	isOpend := C.VideoWriter_IsOpened(vw.p)
	return isOpend != 0
}

// Write the next video frame from the Mat image to the open VideoWriter.
//
// For further details, please see:
// http://docs.opencv.org/master/dd/d9e/classcv_1_1VideoWriter.html#a3115b679d612a6a0b5864a0c88ed4b39
//
func (vw *VideoWriter) Write(img Mat) error {
	vw.mu.Lock()
	defer vw.mu.Unlock()
	C.VideoWriter_Write(vw.p, img.p)
	return nil
}

// OpenVideoCapture return VideoCapture specified by device ID if v is a
// number. Return VideoCapture created from video file, URL, or GStreamer
// pipeline if v is a string.
func OpenVideoCapture(v interface{}) (*VideoCapture, error) {
	switch vv := v.(type) {
	case int:
		return VideoCaptureDevice(vv)
	case string:
		id, err := strconv.Atoi(vv)
		if err == nil {
			return VideoCaptureDevice(id)
		}
		return VideoCaptureFile(vv)
	default:
		return nil, errors.New("argument must be int or string")
	}
}
