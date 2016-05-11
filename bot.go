package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Yamashou/MyClassSearch"
	"github.com/Yamashou/MyStudyRoomSearch"
	"github.com/Yamashou/RandomWord"
	"github.com/Yamashou/SearchFreeRoom"
	"github.com/kurouw/infoSub"
	"github.com/kurouw/reqCafe"
	"github.com/m2mtu/facebookbot/endpoints"
)

// DistributeMenu express functions of bot
type DistributeMenu struct {
	Judgment []string
	Jf       bool
}

func main() {
	os.Setenv("HTTP_PROXY", os.Getenv("FIXIE_URL"))
	os.Setenv("HTTPS_PROXY", os.Getenv("FIXIE_URL"))
	fmt.Println("starting...")
	endpoints.Listen(handleReceiveMessage)
}

func handleReceiveMessage(receivedEvent endpoints.Event) {
	sentEvent := endpoints.Event{}
	sentEvent.SenderID = receivedEvent.RecepientID
	sentEvent.RecepientID = receivedEvent.SenderID
	switch content := receivedEvent.Content.(type) {
	case endpoints.TextContent:
		sentEvent.Content = endpoints.TextContent{Text: getMessageText(content.Text)}
	}
	endpoints.Send(sentEvent)
}

func selectMenu(txt string) string {
	foods := new(DistributeMenu)
	foods.Judgment = []string{"kondate", "こんだて", "献立", "学食", "めにゅー", "メニュー"}
	foods.Jf = false

	tandai := new(DistributeMenu)
	tandai.Judgment = []string{"tandai", "短大", "たんだい"}
	tandai.Jf = false

	computers := new(DistributeMenu)
	computers.Judgment = []string{"演習室", "パソコン", "pc"}
	computers.Jf = false

	eves := new(DistributeMenu)
	eves.Judgment = []string{"hoge"}
	eves.Jf = false

	rooms := new(DistributeMenu)
	rooms.Judgment = []string{"std1", "std2", "std3", "std4", "std5", "std6", "hdw1", "hdw2", "hdw3", "hdw4", "CALL1", "CALL2", "iLab1", "iLab2"}
	rooms.Jf = false

	frooms := new(DistributeMenu)
	frooms.Judgment = []string{"1限", "2限", "3限", "4限", "5限", "6限"}
	frooms.Jf = false

	stringnames := []string{"foods", "tandai", "computers", "eves", "rooms", "frooms"}
	allEvents := []DistributeMenu{*foods, *tandai, *computers, *eves, *rooms, *frooms}

	for i := range allEvents {
		for j := 0; j < len(allEvents[i].Judgment); j++ {
			r := regexp.MustCompile(allEvents[i].Judgment[j])
			if r.MatchString(txt) {
				allEvents[i].Jf = true
			}
		}
	}
	flag := false
	for i := range allEvents {
		if allEvents[i].Jf {
			allEvents[i].Jf = false
			flag = true
			return stringnames[i]
		}
	}
	if !flag {
		cflag := false
		name := txt
		name = string([]rune(name)[:1])
		if name == "s" || name == "m" {
			cflag = true
			return "classes"
		}
		if !cflag {
			return "Subject!"
		}
	}
	return "n"
}

func getMessageText(receivedText string) string {
	var sub string

	selectRes := selectMenu(receivedText)
	fmt.Println("selected: " + selectRes)
	if selectRes == "foods" {
		var res []string
		res = reqCafe.RtCafeInfo(time.Now())

		b := make([]byte, 0, 30)
		for v := 0; v < len(res); v++ {
			b = append(b, res[v]...)
			b = append(b, '\n')
		}
		return string(b)

	} else if selectRes == "tandai" {
		var res []string
		res = reqCafe.RtTnCafeInfo(time.Now())

		b := make([]byte, 0, 30)
		for v := 0; v < len(res); v++ {
			b = append(b, res[v]...)
			b = append(b, '\n')
		}
		return string(b)

	} else if selectRes == "rooms" {
		room := MyStudyRoomSearch.RtRoom(receivedText)
		b := make([]byte, 0, 30)
		for v := 0; v < len(room); v++ {
			b = append(b, strconv.Itoa(v+1)+"限: "...)
			b = append(b, room[v]...)
			b = append(b, '\n')
		}
		return string(b)
	} else if selectRes == "frooms" {

		var frooms [15]string
		var num int
		name := receivedText
		name = string([]rune(name)[:1])
		num, _ = strconv.Atoi(name)
		frooms = SearchFreeRoom.Serect(num)

		b := make([]byte, 0, 30)
		for v := 0; v < len(frooms); v++ {
			b = append(b, frooms[v]...)
			b = append(b, '\n')
		}
		return string(b)

	}

	if selectRes == "Subject!" {
		sub = infoSub.ReturnSubInfo(receivedText)
	}

	if selectRes == "classes" {

		stdClass := MyClassSearch.RtClass(receivedText)

		b := make([]byte, 0, 30)
		for v := 0; v < len(stdClass); v++ {
			b = append(b, strconv.Itoa(v+1)+"限: "...)
			b = append(b, stdClass[v]...)
			b = append(b, '\n')
		}
		return string(b)

	}
	if sub != receivedText {
		return sub
	} else {
		return RandomWord.ReturnWord(receivedText)
	}
}
