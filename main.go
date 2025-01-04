package main

import (
	"fmn/journalbot/telegram"
	"fmn/journalbot/updater"
)

func main() {
	updater := updater.NewUpdater()
	telegram := telegram.InitService(updater)
	telegram.Start()
}
