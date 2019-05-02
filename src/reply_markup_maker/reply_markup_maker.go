package reply_markup_maker

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func MakeReplyMarkupSmart(keyword string, buttonsLen int, index int) tgbotapi.InlineKeyboardMarkup {
	var startPos int
	var endPos int
	var buttons []tgbotapi.InlineKeyboardButton
	for i := 0; i < buttonsLen; i++ {
		callbackData := keyword + "," + strconv.Itoa(i) //probleem,1
		but := tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i), callbackData)
		buttons = append(buttons, but)
	}
	buttons[index].Text = ">" + buttons[index].Text + "<"
	if buttonsLen <= 5 {
		startPos = 0
		endPos = buttonsLen
	} else if index-3 >= 0 && index + 1 < buttonsLen {
		startPos = index - 3
		endPos = index + 2
	} else if index-2 <= 0 {
		startPos = 0
		endPos = 5
	} else if index+2 > buttonsLen {
		startPos = buttonsLen - 5
		endPos = buttonsLen
	}

	buttons = buttons[startPos:endPos]
	if startPos > 0 {
		buttons[0].Text = "<<" + buttons[0].Text
	}

	if endPos < buttonsLen {
		buttons[len(buttons)-1].Text += ">>"
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}
