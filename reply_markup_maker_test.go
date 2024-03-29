package main

import (
	"testing"
)

func TestMakeReplyMarkup_len5(t *testing.T) {
	articlesLen := 5
	kbd := MakeReplyMarkup("abc", articlesLen, 4)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}
}

func TestMakeReplyMarkup_len7_index1(t *testing.T) {
	articlesLen := 7
	kbd := MakeReplyMarkup("abc", articlesLen, 1)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}
	if kbd.InlineKeyboard[0][0].Text != "0" && kbd.InlineKeyboard[0][4].Text != "4>>" {
		t.Fail()
	}

}

func TestMakeReplyMarkup_len7_index6(t *testing.T) {
	articlesLen := 7
	kbd := MakeReplyMarkup("abc", articlesLen, 6)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "<<2" && kbd.InlineKeyboard[0][4].Text != "6" {
		t.Fail()
	}
}

func TestMakeReplyMarkup_len20_index18(t *testing.T) {
	articlesLen := 20
	kbd := MakeReplyMarkup("abc", articlesLen, 18)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "<<15" && kbd.InlineKeyboard[0][4].Text != "19" {
		t.Fail()
	}
}

func TestMakeReplyMarkup_len20_index8(t *testing.T) {
	articlesLen := 20
	kbd := MakeReplyMarkup("abc", articlesLen, 8)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "<<5" && kbd.InlineKeyboard[0][4].Text != "9>>" {
		t.Fail()
	}
}

func TestMakeReplyMarkup_len20_index2(t *testing.T) {
	articlesLen := 20
	kbd := MakeReplyMarkup("abc", articlesLen, 2)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "0" && kbd.InlineKeyboard[0][4].Text != "4>>" {
		t.Fail()
	}
}

func TestMakeReplyMarkup_current_elem_empty_cb(t *testing.T) {
	articlesLen := 5
	selectedElemIndex := 2
	kbd := MakeReplyMarkup("abc", articlesLen, selectedElemIndex)
	if *kbd.InlineKeyboard[0][selectedElemIndex].CallbackData != "none" {
		t.Fail()
	}
}
