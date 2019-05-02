package main

import (
	"testing"
	"strconv"
)

func TestMakeReplyMarkupNice_len4_start0(t *testing.T) {
	articlesLen := 4
	kbd := MakeReplyMarkupNice("abc", articlesLen, 0)
	if len(kbd.InlineKeyboard[0]) != articlesLen {
		t.Fail()
	}

	for i := 0; i < articlesLen; i++ {
		if kbd.InlineKeyboard[0][i].Text != strconv.Itoa(i) {
			t.Fail()
		}
	}
}

func TestMakeReplyMarkupNice_len6_start0(t *testing.T) {

	articlesLen := 6
	kbd := MakeReplyMarkupNice("abc", articlesLen, 0)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "0" {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][4].Text != "4>>" {
		t.Fail()
	}
}

func TestMakeReplyMarkupNice_len6_start5(t *testing.T) {
	//[<<1][2][3][4][5]
	articlesLen := 6
	kbd := MakeReplyMarkupNice("abc", articlesLen, 5)
	if len(kbd.InlineKeyboard[0]) != 5 {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][0].Text != "<<1" {
		t.Fail()
	}

	if kbd.InlineKeyboard[0][4].Text != "5" {
		t.Fail()
	}
}
