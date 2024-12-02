# Dependencies

To build this application the following dependencies are required:

```shell
$ go get gocv.io/x/gocv
```

GoCV needs to be compiled first (change the version number to the one that is installed):

```shell
$ cd $GOPATH/pkg/mod/gocv.io/x/gocv@v0.39.0
$ make install
```

# Usage

To run the application in `yuyv` mode with `exposure` 2048:

```shell
$ task build
$ ./bin/cameratest -mode=mjpg -exposure=2048
```

## Options
* `-mode=<value>` - Output mode. Options are `yuyv` and `mjpg`. No default, must be set.
* `-exposure=<value>` - Exposure value (`20` to `10000`). Default is `512`. Recommend one of `20`, `512`, `1024`, `2048`, `4096`, or `8192`.
* `-delay=<value>` - Delay to allow camera to initialize after changing properties (`0` to `10000` milliseconds). Default is `0`.
* `-vid=<value>` - Camera USB Vendor ID (in hexadecimal). Default is `1bcf` (Fangtec 2nd generation camera).
* `-pid=<value>` - Camera USB Product ID (in hexadecimal). Default is `0b09` (Fangtec 2nd generation camera).
* `-fix` - If the exposure is not set correctly, this will attempt to fix it.

# Errors

Some example errors that have been observed with the current Fangtec 2nd generation camera firmware are shown below.
This is not an exhaustive list, and other errors may occur until the firmware is fixed. Please report any errors not
already listed below to `trent@assaya.com`, and include the log output and the camera firmware version.

### Select Timeout

If the select() timeout error occurs, it will look like this:

```text
2024/11/04 14:41:06 Using camera device: /dev/video4
2024/11/04 14:41:06 Identified CAM: /dev/video4
2024/11/04 14:41:06 Initialized CAM: MJPG 1280x720 30fps
Gtk-Message: 14:41:06.756: Failed to load module "canberra-gtk-module"
2024/11/04 14:41:06 Press ESC to exit.
[ WARN:0@10.539] global cap_v4l.cpp:1136 tryIoctl VIDEOIO(V4L2:/dev/video4): select() timeout.
2024/11/04 14:41:16 ERROR reading frame: 1
[ WARN:0@20.549] global cap_v4l.cpp:1136 tryIoctl VIDEOIO(V4L2:/dev/video4): select() timeout.
2024/11/04 14:41:26 ERROR reading frame: 2
[ WARN:0@30.559] global cap_v4l.cpp:1136 tryIoctl VIDEOIO(V4L2:/dev/video4): select() timeout.
2024/11/04 14:41:36 ERROR reading frame: 3
[ WARN:0@40.569] global cap_v4l.cpp:1136 tryIoctl VIDEOIO(V4L2:/dev/video4): select() timeout.
2024/11/04 14:41:46 ERROR reading frame: 4
[ WARN:0@50.579] global cap_v4l.cpp:1136 tryIoctl VIDEOIO(V4L2:/dev/video4): select() timeout.
2024/11/04 14:41:56 ERROR reading frame: 5
```

### Set Control Value

If a set control error occurs it will look like this:

```text
2024/11/04 15:18:44 Using camera device: /dev/video4
2024/11/04 15:18:44 SetCamSetting Error: White Balance Temperature - device: /dev/video4: set control value: id 9963802: broken pipe
2024/11/04 15:18:44 Identified CAM: /dev/video4
2024/11/04 15:18:44 Initialized CAM: MJPG 1280x720 30fps
Gtk-Message: 15:18:44.453: Failed to load module "canberra-gtk-module"
2024/11/04 15:18:44 Press ESC to exit.
```

### Corrupt JPEG Data

If corrupt JPEG data occurs it will look like this:

```text
2024/11/04 15:47:22 Using camera device: /dev/video4
2024/11/04 15:47:22 Identified CAM: /dev/video4
2024/11/04 15:47:22 Initialized CAM: MJPG 1280x720 30fps
2024/11/04 15:47:22 Camera controls:
2024/11/04 15:47:22 Brightness: 0, Expected: 0
2024/11/04 15:47:22 Contrast: 15, Expected: 15
2024/11/04 15:47:22 Gamma: 100, Expected: 100
2024/11/04 15:47:22 Hue: 0, Expected: 0
2024/11/04 15:47:22 Saturation: 65, Expected: 65
2024/11/04 15:47:22 Sharpness: 30, Expected: 30
2024/11/04 15:47:22 Temperature: 6000, Expected: 6000
2024/11/04 15:47:22 Exposure: 512, Expected: 512
Gtk-Message: 15:47:22.610: Failed to load module "canberra-gtk-module"
2024/11/04 15:47:22 Press ESC to exit.
2024/11/04 15:47:30 ERROR reading frame: 244
2024/11/04 15:47:44 ERROR reading frame: 661
Corrupt JPEG data: premature end of data segment
```

# Packaging

To package the application for distribution on Ubuntu 22.04 LTS (Jammy Jellyfish) or Ubuntu 24.04 LTS (Noble Numbat)
run the following command:

```shell
$ task package
```

Only Ubuntu 22.04 LTS (Jammy Jellyfish) and Ubuntu 24.04 LTS (Noble Numbat) are supported at this time.