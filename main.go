package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)
const (
    maxFilenameLength = 255
    filePermission    = 0644
    thumbnailFileExt  = ".jpg"
    shortsFormat      = "https://www.youtube.com/shorts/"
    videoFormat       = "https://www.youtube.com/watch?"
    imageDownloadLink = "https://img.youtube.com/vi"
    resolution        = "maxresdefault.jpg"
)
func downloadThumbnail(videoID string) error {
    url := fmt.Sprintf("%s/%s/%s", imageDownloadLink, videoID, resolution)
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    filename := videoID + thumbnailFileExt
    if len(filename) > maxFilenameLength {
        filename = filename[:maxFilenameLength]
    }
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    if err != nil {
        return err
    }

    return nil
}
func getUserVideoURL() (string, error) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter YouTube video URL to download thumbnail from: ")
    videoURL, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }

    var videoID string

    if strings.HasPrefix(videoURL, shortsFormat) {
        splitURL := strings.Split(videoURL, "/")
        videoID = splitURL[len(splitURL)-1]
    } else if strings.HasPrefix(videoURL, videoFormat) {
        valueIndex := strings.Index(videoURL, "v=")
        if valueIndex != -1 {
            value := videoURL[valueIndex+2:]
            nextAmpersand := strings.Index(value, "&")
            if nextAmpersand != -1 {
                value = value[:nextAmpersand]
            }
            videoID = value
        }
    }

    videoID = strings.TrimSpace(videoID)
    if videoID == "" {
        return "", fmt.Errorf("Invalid YouTube URL, make sure your inputs are in the format %s<video_id> or %s<video_id>", shortsFormat, videoFormat)
    }

    return videoID, nil
}
func main() {

    videoID, err := getUserVideoURL()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Printf("Downloading thumbnail for video ID %s now...\n", videoID)
    if err := downloadThumbnail(videoID); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println("Thumbnail downloaded successfully!")
}