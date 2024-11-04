package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"gocv.io/x/gocv"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type ControlInfo struct {
	key     string
	name    string
	display bool
}

type FourCC struct {
	mjpg int
	yuyv int
}

// Replace these with your specific vendor and model IDs
const (
	vendorID = "1bcf" // Example vendor ID for Logitech
	modelID  = "0b09" // Example model ID for a Logitech camera
)

// Original camera
//const (
//	vendorID = "058f" // Example vendor ID for Logitech
//	modelID  = "3822" // Example model ID for a Logitech camera
//)

var cameraMode = FourCC{
	mjpg: 1196444237,
	yuyv: 2020216696,
}

var controlNames = map[v4l2.CtrlID]ControlInfo{
	// Camera Controls
	v4l2.CtrlCameraExposureAuto:         ControlInfo{"auto_exposure", "Auto Exposure", false},
	v4l2.CtrlCameraExposureAutoPriority: ControlInfo{"exposure_dynamic_framerate", "Exposure Dynamic Framerate", false},
	v4l2.CtrlCameraExposureAbsolute:     ControlInfo{"exposure", "Exposure Absolute", true},

	// User Controls
	v4l2.CtrlAutoWhiteBalance:        ControlInfo{"auto_white_balance", "Auto White Balance", false},
	v4l2.CtrlBrightness:              ControlInfo{"brightness", "Brightness", true},
	v4l2.CtrlContrast:                ControlInfo{"contrast", "Contrast", true},
	v4l2.CtrlGamma:                   ControlInfo{"gamma", "Gamma", true},
	v4l2.CtrlHue:                     ControlInfo{"hue", "Hue", true},
	v4l2.CtrlSaturation:              ControlInfo{"saturation", "Saturation", true},
	v4l2.CtrlSharpness:               ControlInfo{"sharpness", "Sharpness", true},
	v4l2.CtrlWhiteBalanceTemperature: ControlInfo{"whitebalancetemperature", "White Balance Temperature", true},
}

var controlValues = map[string]int{
	"brightness":  0,
	"contrast":    15,
	"gamma":       100,
	"hue":         0,
	"saturation":  65,
	"sharpness":   30,
	"temperature": 6000,
}

var exposureValue int
var fixValue bool

func main() {
	// Define and parse the command-line flags for mode and exposure
	modeFlag := flag.String("mode", "mjpg", "Mode: mjpg or yuyv")
	exposureFlag := flag.Int("exposure", 512, "Exposure value")
	fixFlag := flag.Bool("fix", false, "Fix the exposure value")
	flag.Parse()

	// Determine the camera mode based on the flag
	var modeValue int

	switch *modeFlag {
	case "mjpg":
		modeValue = cameraMode.mjpg
	case "yuyv":
		modeValue = cameraMode.yuyv
	default:
		log.Fatalf("Invalid mode: %s. Please choose 'mjpg' or 'yuyv'.", *modeFlag)
	}

	switch {
	case *exposureFlag < 20:
		log.Fatalf("Exposure value must be greater than 20. Please choose a value greater than 20.")
	case *exposureFlag > 10000:
		log.Fatalf("Exposure value must be less than 8192. Please choose a value less than 8192.")
	default:
		exposureValue = *exposureFlag
	}

	switch *fixFlag {
	case true:
		fixValue = true
	default:
		fixValue = false
	}

	// Find the camera device path
	cameraDevice, err := devicePath(vendorID, modelID)
	if err != nil {
		log.Println("Error:", err.Error())
		return
	}

	log.Println("Using camera device:", cameraDevice)

	// Initialize the camera device
	err = deviceInit(cameraDevice)
	if err != nil {
		log.Println("Error:", err.Error())
		return
	}

	// Open a video capture device (0 is usually the default USB camera)
	webcam, err := gocv.OpenVideoCapture(cameraDevice)
	if err != nil {
		log.Println("Error opening video capture device:", err.Error())
		return
	}
	defer webcam.Close()

	// Set camera properties (adjust according to your needs)
	webcam.Set(gocv.VideoCaptureFOURCC, float64(modeValue))
	webcam.Set(gocv.VideoCaptureFrameWidth, 1280)
	webcam.Set(gocv.VideoCaptureFrameHeight, 720)

	actualCodec := webcam.CodecString()
	actualWidth := int(webcam.Get(gocv.VideoCaptureFrameWidth))
	actualHeight := int(webcam.Get(gocv.VideoCaptureFrameHeight))
	actualFPS := int(webcam.Get(gocv.VideoCaptureFPS))

	log.Println(fmt.Sprintf("Identified CAM: %v", cameraDevice))
	log.Println(fmt.Sprintf("Initialized CAM: %s %dx%d %dfps", actualCodec, actualWidth, actualHeight, actualFPS))

	// Report on the actual state of the device
	err = deviceState(cameraDevice)
	if err != nil {
		log.Println("Error:", err.Error())
		return
	}

	// Prepare a Mat to store frames
	frame := gocv.NewMat()
	defer frame.Close()

	// Create a window to display the video feed
	window := gocv.NewWindow("USB Camera Feed")
	defer window.Close()
	window.ResizeWindow(actualWidth, actualHeight)

	log.Println("Press ESC to exit.")

	var frameCount int64

	// Main loop to capture frames
	for {
		frameCount++

		// Read a frame from the camera
		if ok := webcam.Read(&frame); !ok {
			log.Println("ERROR reading frame:", frameCount)
			continue
		}

		if frame.Empty() {
			log.Println("Empty frame:", frameCount)
			continue
		}

		//// Convert the frame to BGR color format
		//gocv.CvtColor(frame, &frame, gocv.ColorBGRToRGB)

		// Display the frame in the window
		window.IMShow(frame)

		// Wait for key press
		key := window.WaitKey(1)
		if key == 27 { // ESC to exit
			break
		}
	}

	log.Println("Exited.")
}

// devicePath finds the camera device path based on vendor and model IDs
func devicePath(vendorID, modelID string) (string, error) {
	// List all video devices in /dev
	videoDevices, err := filepath.Glob("/dev/video*")
	if err != nil {
		return "", fmt.Errorf("error listing video devices: %v", err)
	}

	// Regular expressions to match vendor and model ID
	vendorRegex := regexp.MustCompile(fmt.Sprintf("ID_VENDOR_ID=%s", vendorID))
	modelRegex := regexp.MustCompile(fmt.Sprintf("ID_MODEL_ID=%s", modelID))

	// Check each video device with udevadm
	for _, videoDevice := range videoDevices {
		cmd := exec.Command("udevadm", "info", "--query=all", "--name="+videoDevice)
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Error running udevadm on %s: %v\n", videoDevice, err)
			continue
		}

		// Scan over the output to
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		var vendorFound bool
		var modelFound bool

		for scanner.Scan() {
			line := scanner.Text()
			if vendorRegex.MatchString(line) {
				vendorFound = true
			}
			if modelRegex.MatchString(line) {
				modelFound = true
			}

			if vendorFound && modelFound {
				return videoDevice, nil
			}
		}
	}
	return "", fmt.Errorf("camera with vendor ID %s and model ID %s not found", vendorID, modelID)
}

func deviceInit(devicePath string) error {
	videoDevice, err := device.Open(devicePath)
	if err != nil {
		log.Println("Error opening video device:", err.Error())
		return err
	}
	defer videoDevice.Close()

	// Camera controls
	setControl(videoDevice, v4l2.CtrlCameraExposureAuto, 1)         // off
	setControl(videoDevice, v4l2.CtrlCameraExposureAutoPriority, 0) // off

	// User controls
	setControl(videoDevice, v4l2.CtrlAutoWhiteBalance, 0) // off
	setControl(videoDevice, v4l2.CtrlBrightness, controlValues["brightness"])
	setControl(videoDevice, v4l2.CtrlContrast, controlValues["contrast"])
	setControl(videoDevice, v4l2.CtrlGamma, controlValues["gamma"])
	setControl(videoDevice, v4l2.CtrlHue, controlValues["hue"])
	setControl(videoDevice, v4l2.CtrlSaturation, controlValues["saturation"])
	setControl(videoDevice, v4l2.CtrlSharpness, controlValues["sharpness"])
	setControl(videoDevice, v4l2.CtrlWhiteBalanceTemperature, controlValues["temperature"])

	// For some reason, we need to set the exposure time later on...
	setControl(videoDevice, v4l2.CtrlCameraExposureAbsolute, exposureValue)

	return nil
}

func deviceState(devicePath string) error {
	videoDevice, err := device.Open(devicePath)
	if err != nil {
		log.Println("Error opening video device:", err.Error())
		return err
	}
	defer videoDevice.Close()

	if getControl(videoDevice, v4l2.CtrlBrightness) != controlValues["brightness"] {
		logMessage("Brightness", controlValues["brightness"], getControl(videoDevice, v4l2.CtrlBrightness))
	}

	if getControl(videoDevice, v4l2.CtrlContrast) != controlValues["contrast"] {
		logMessage("Contrast", controlValues["contrast"], getControl(videoDevice, v4l2.CtrlContrast))
	}

	if getControl(videoDevice, v4l2.CtrlGamma) != controlValues["gamma"] {
		logMessage("Gamma", controlValues["gamma"], getControl(videoDevice, v4l2.CtrlGamma))
	}

	if getControl(videoDevice, v4l2.CtrlHue) != controlValues["hue"] {
		logMessage("Hue", controlValues["hue"], getControl(videoDevice, v4l2.CtrlHue))
	}

	if getControl(videoDevice, v4l2.CtrlSaturation) != controlValues["saturation"] {
		logMessage("Saturation", controlValues["saturation"], getControl(videoDevice, v4l2.CtrlSaturation))
	}

	if getControl(videoDevice, v4l2.CtrlSharpness) != controlValues["sharpness"] {
		logMessage("Sharpness", controlValues["sharpness"], getControl(videoDevice, v4l2.CtrlSharpness))
	}

	if getControl(videoDevice, v4l2.CtrlWhiteBalanceTemperature) != controlValues["temperature"] {
		logMessage("Temperature", controlValues["temperature"], getControl(videoDevice, v4l2.CtrlWhiteBalanceTemperature))
	}

	if getControl(videoDevice, v4l2.CtrlCameraExposureAbsolute) != exposureValue {
		logMessage("Exposure", exposureValue, getControl(videoDevice, v4l2.CtrlCameraExposureAbsolute))

		if fixValue {
			var adjustedExposure int

			if exposureValue < 512 {
				adjustedExposure = 512

			} else if exposureValue > 4096 {
				adjustedExposure = 8192
			} else {
				adjustedExposure = exposureValue * 2
			}

			log.Println("Adjusting exposure to", adjustedExposure)
			setControl(videoDevice, v4l2.CtrlCameraExposureAbsolute, adjustedExposure)

			log.Println("Adjusting exposure back to", exposureValue)
			setControl(videoDevice, v4l2.CtrlCameraExposureAbsolute, exposureValue)

			if getControl(videoDevice, v4l2.CtrlCameraExposureAbsolute) != exposureValue {
				logMessage("Exposure", exposureValue, getControl(videoDevice, v4l2.CtrlCameraExposureAbsolute))
			} else {
				fmt.Println("Exposure adjusted successfully.")
			}
		}
	}

	return nil
}

func logMessage(name string, expected int, actual int) {
	log.Printf("ERROR: %s is not set correctly. Set value: %d, Actual value: %d", name, expected, actual)
}

func setControl(videoDevice *device.Device, controlID v4l2.CtrlID, value int) {
	err := videoDevice.SetControlValue(controlID, v4l2.CtrlValue(value))
	if err != nil {
		log.Println("SetCamSetting Error:", controlNames[controlID].name, "-", err)
	}
}

func getControl(videoDevice *device.Device, controlID v4l2.CtrlID) int {
	control, err := videoDevice.GetControl(controlID)
	if err != nil {
		log.Println("GetCamSetting Error:", controlNames[controlID].name, "-", err)
		return 0
	}

	return int(control.Value)
}
