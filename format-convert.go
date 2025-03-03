package main

import (
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os/exec"
	"runtime"
)

func (a *App) Checkffmpeg() bool {
	return Checkffmpeg()
}

// 检查系统中是否安装 ffmpeg
// （临时方案）
func Checkffmpeg() bool {
	switch runtime.GOOS {
	case "windows":
		return checkffmpegOnWindows()
	case "darwin":
		return checkffmpegOnMacOS()
	default:
		return false
	}
}

// windows
func checkffmpegOnWindows() bool {
	cmd := exec.Command("where", "ffmpeg")
	setHideWindow(cmd)
	_, err := cmd.Output()
	return err == nil
}

// MacOS
func checkffmpegOnMacOS() bool {
	cmd := exec.Command("which", "ffmpeg")
	_, err := cmd.Output()
	return err == nil
}

// 调用 ffmpeg 转码
func ConventFile(inputFile, outputFile string) error {
	stream := ffmpeg.Input(inputFile).Output(outputFile, ffmpeg.KwArgs{"qscale": "0"})
	cmd := stream.Compile()
	setHideWindow(cmd)
	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}

// 处理flac metadata (清除metadata)
func HandleFlacMetadata(inputFile string) error {
	// mv inputFile tmp.flac
	// ffmpeg -i tmp.flac -c:a flac -map_metadata -1 -y outputFile
	tmpFile := inputFile + ".tmp.flac"
	err := exec.Command("mv", inputFile, tmpFile).Run()
	if err != nil {	// 如果移动失败
		return err
	}
	stream := ffmpeg.Input(tmpFile).Output(inputFile, ffmpeg.KwArgs{"c:a": "flac", "map_metadata": "-1"})
	cmd := stream.Compile()
	setHideWindow(cmd)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
