package main

import (
	"animebot/downloadhelper"
	"animebot/qqapi"
	"flag"
	"fmt"
	"log"
	"os"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *botApi.BotAPI

func main() {

	str := flag.String("token", os.Getenv("BOT_TOKEN"), "Telegram bot token")
	flag.Parse()
	token := *str

	initializeBot(token)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := botApi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.Photo != nil { // ignore any non-image Messages
			HandleImage(update.Message)
		} else {
			SendReply(update.Message.Chat.ID, "Send me ONE image pls")
		}
	}
}

func HandleImage(message *botApi.Message) {

	imageSizes := message.Photo

	chatId := message.Chat.ID

	fileId := imageSizes[len(imageSizes)-1].FileID

	file, err := bot.GetFileDirectURL(fileId)
	if err != nil {
		SendReply(chatId, err.Error())
		return
	}

	log.Println("Got file url:" + file)
	SendReply(chatId, "Processing file, please wait...")

	base64 := downloadhelper.DownloadImageAsBase64(file)
	results := qqapi.RequestImage(base64)

	log.Println("Saving media...")
	SendReply(chatId, "Saving media...")

	arr := append(addImages(results.ImgUrls), addVideos(results.VideoUrls)...)

	log.Println("Creating album...")
	SendReply(chatId, "Creating album...")

	album := botApi.NewMediaGroup(chatId, arr)
	if _, err := bot.Send(album); err != nil {
		log.Print(err)
	}
}

func addImages(imgUrls []string) []interface{} {
	arr := []interface{}{}

	for i, img := range imgUrls {
		file := downloadhelper.DownloadFile(img)
		fileBytes := botApi.FileBytes{
			Name:  fmt.Sprintf("image%d", i),
			Bytes: file,
		}

		image := botApi.NewInputMediaPhoto(fileBytes)

		arr = append(arr, image)
	}

	return arr
}

func addVideos(vidUrls []string) []interface{} {
	arr := []interface{}{}

	for i, vid := range vidUrls {
		file := downloadhelper.DownloadFile(vid)
		fileBytes := botApi.FileBytes{
			Name:  fmt.Sprintf("video%d", i),
			Bytes: file,
		}
		image := botApi.NewInputMediaVideo(fileBytes)
		arr = append(arr, image)
	}

	return arr
}

func SendReply(chatId int64, text string) {
	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	reply := botApi.NewMessage(chatId, text)
	if _, err := bot.Send(reply); err != nil {
		log.Print(err)
	}
}

func initializeBot(token string) {
	//token := os.Getenv("BOT_TOKEN")
	BotAPI, err := botApi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	bot = BotAPI
}
