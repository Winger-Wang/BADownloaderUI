package main

import (
	"bytes"
	"html/template"
	"os"
	"path"
	// "path/filepath"
	"strconv"
	"time"
	"os/exec"

	tag "github.com/gcottom/audiometa"
	// tag_flac "github.com/gcottom/flacmeta"
	// tag_mp3 "github.com/gcottom/mp3meta"
)

// 修改 TAG
func ChangeTag(cfg *Config, opt *DownloadOption, v *VideoInformation) error {

	// 准备参数
	file := cfg.FileConfig.CachePath + "/music/" + strconv.Itoa(v.Cid) + v.Format
	songCover := cfg.FileConfig.CachePath + "/cover/" + strconv.Itoa(v.Cid) + ".jpg"
	songName := v.Meta.SongName
	songAuthor := v.Meta.Author
	// turn publishTime (string) to the publishTime (int64)
	publishTime, _ := strconv.ParseInt(v.PublishTime, 10, 64)

	// 打开歌曲元数据
	tags, err := tag.OpenTag(file)
	if err != nil {
		return err
	}

	// 封面
	if opt.SongCover {
		err := tags.SetAlbumArtFromFilePath(songCover)
		if err != nil {
			return err
		}
	}
	// 歌曲名
	if opt.SongName {
		tags.SetTitle(songName)
	}
	// 艺术家
	if opt.SongAuthor {
		tags.SetArtist(songAuthor)
	}

	// 这里开始是我添加的代码，用于修改歌曲的发布时间和专辑名
	// 专辑名
	tags.SetAlbum(songName)
	publishDate := time.Unix(publishTime, 0).Format("2006-01-02")
	publishYear := time.Unix(publishTime, 0).Year()
	
	// TODO: 将歌曲 tag 数据整理为结构体
	// TODO: 修改作词人，作曲人等，以及自动适配

	// 保存更改
	err = tags.Save()
	if err != nil {
		return err
	}

	// use ffmpeg to change metadata of release time
	mv := exec.Command("mv", file, file + ".temp." + v.Format)
	err = mv.Run()
	if err != nil {
		return err
	}
    cmd := exec.Command("ffmpeg", "-i", file+".temp."+v.Format, "-metadata", "date="+publishDate, "-metadata", "year="+strconv.Itoa(publishYear), "-y", file)
	err = cmd.Run()
	if err != nil {
		return err
	}

	// // 发布时间
	// // turn publishTime (timeStamp) to the publishYear
	// publishDate := time.Unix(publishTime, 0).Format("2006-01-02")
	// publishYear := time.Unix(publishTime, 0).Year()
    // // tags.SetYear(strconv.Itoa(publishYear))
	// // mp3
	// if v.Format == AudioType.mp3 {
	// 	// open the mp3 file using the mp3meta decoder
	// 	mp3file, _ := os.Open(file)
	// 	tags, err := tag_mp3.ParseMP3(mp3file)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	tags.SetYear(publishYear)
	// 	tag_mp3.SaveMP3(tags, mp3file)
	// 	mp3file.Close()
	// }
	// // flac
	// if v.Format == AudioType.flac {
	// 	// open the flac file using the flacmeta decoder
	// 	flacfile, _ := os.Open(file)
	// 	tags, err := tag_flac.ReadFLAC(flacfile)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	tags.SetTitle("Test title")
	// 	tags.SetDate(publishDate)
	// 	tags.Save(flacfile)
	// 	flacfile.Close()
	// }

	return nil
}

type FileName struct {
	Title    string
	Subtitle string
	Quality  string
	ID       int
	Format   string
}

// 输出文件
func OutputFile(cfg *Config, v *VideoInformation, fileName FileName) error {
	// 处理模板和生成文件名
	tmpl, err := template.New("filename").Parse(cfg.FileConfig.FileNameTemplate)
	if err != nil {
		return err
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, fileName)
	if err != nil {
		return err
	}

	// 添加路径
	sourcePath := path.Join(cfg.FileConfig.CachePath, "music", strconv.Itoa(v.Cid)+v.Format)
	destPath := path.Join(cfg.FileConfig.DownloadPath, output.String())

	// 重命名歌曲文件并移动位置
	err = os.Rename(sourcePath, destPath)
	if err != nil {
		return err
	}
	return nil
}

// 修改 TAG
func SingleChangeTag(cfg *Config, opt *DownloadOption, auid, songName, songAuthor string) error {

	// 准备参数
	file := cfg.FileConfig.CachePath + "/single/music/" + auid + AudioType.m4a
	songCover := cfg.FileConfig.CachePath + "/single/cover/" + auid + ".jpg"

	// 打开歌曲元数据
	tags, err := tag.OpenTag(file)
	if err != nil {
		return err
	}

	// 封面
	if opt.SongCover {
		tags.SetAlbumArtFromFilePath(songCover)
	}
	// 歌曲名
	if opt.SongName {
		tags.SetTitle(songName)
	}
	// 艺术家
	if opt.SongAuthor {
		tags.SetArtist(songAuthor)
	}

	// TODO: 将歌曲 tag 数据整理为结构体
	// TODO: 修改作词人，作曲人等，以及自动适配

	// 保存更改
	err = tags.Save()
	if err != nil {
		return err
	}

	return nil
}

// 输出文件
func SingleOutputFile(cfg *Config, uuid, Title string) error {

	sourcePath := path.Join(cfg.FileConfig.CachePath, "single/music", uuid+AudioType.m4a)
	destPath := path.Join(cfg.FileConfig.DownloadPath, Title+AudioType.mp3)

	// 重命名歌曲文件并移动位置
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		return err
	}
	return nil
}
