package reply_markup_maker

import (
	"testing"
)

func TestMakeReplyMarkupSmart_len5(t *testing.T) {
	articlesLen := 5
	kbd := MakeReplyMarkupSmart("abc", articlesLen, 4)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}
}

func TestMakeReplyMarkupSmart_len7_index1(t *testing.T) {
	articlesLen := 7
	kbd := MakeReplyMarkupSmart("abc", articlesLen, 1)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}
	if kbd.InlineKeyboard[0][0].Text != "0" && kbd.InlineKeyboard[0][4].Text != "4>>" {
		t.Fail()
	}


}

func TestMakeReplyMarkupSmart_len7_index6(t *testing.T) {
	articlesLen := 7
	kbd := MakeReplyMarkupSmart("abc", articlesLen, 6)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "<<2" && kbd.InlineKeyboard[0][4].Text != "6" {
		t.Fail()
	}
}


func TestMakeReplyMarkupSmart_len20_index18(t *testing.T) {
	articlesLen := 20
	kbd := MakeReplyMarkupSmart("abc", articlesLen, 18)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "<<15" && kbd.InlineKeyboard[0][4].Text != "19" {
		t.Fail()
	}
}

func TestMakeReplyMarkupSmart_len20_index8(t *testing.T) {
	articlesLen := 20
	kbd := MakeReplyMarkupSmart("abc", articlesLen, 8)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "<<5" && kbd.InlineKeyboard[0][4].Text != "9>>" {
		t.Fail()
	}
}

func TestMakeReplyMarkupSmart_len20_index2(t *testing.T) {
	articlesLen := 20
	kbd := MakeReplyMarkupSmart("abc", articlesLen, 2)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "0" && kbd.InlineKeyboard[0][4].Text != "4>>" {
		t.Fail()
	}
}

