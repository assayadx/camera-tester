# Usage

To run the application in `yuyv` mode with `exposure` 2048:

```shell
go run main.go -mode=yuyv -exposure=2048

To run the application in `mjpg` mode with `exposure` 2048:

```shell
go run main.go -mode=mjpg -exposure=2048

```
Sometimes the exposure will not get set correctly. To fix this, you can add the `fix` command option:

```shell
go run main.go -mode=yuyv -exposure=2048 -fix
```

# Errors

If the select() timeout error occurs, it will look like this:

```shell
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

If a set control error occurs it will look like this:

```shell
2024/11/04 15:18:44 Using camera device: /dev/video4
2024/11/04 15:18:44 SetCamSetting Error: White Balance Temperature - device: /dev/video4: set control value: id 9963802: broken pipe
2024/11/04 15:18:44 Identified CAM: /dev/video4
2024/11/04 15:18:44 Initialized CAM: MJPG 1280x720 30fps
Gtk-Message: 15:18:44.453: Failed to load module "canberra-gtk-module"
2024/11/04 15:18:44 Press ESC to exit.
```

If corrupt JPEG data occurs it will look like this:

```shell
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