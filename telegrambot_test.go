package main

import (
	"gopkg.in/tucnak/telebot.v2"
	"reflect"
	"testing"
)

func TestDecideWhetherToRespond(t *testing.T) {
	tables := []struct {
		testcase   string
		chattiness float64
		name       string
		m          *telebot.Message
		expected   bool
	}{
		{
			"Feeling chatty",
			1,
			"@bot",
			&telebot.Message{Text: "test"},
			true,
		},
		{
			"Not feeling chatty",
			0,
			"@bot",
			&telebot.Message{
				Text: "test",
				Chat: &telebot.Chat{Type: ""},
			},
			false,
		},
		{
			"Private chat, not feeling chatty",
			0,
			"@bot",
			&telebot.Message{
				Text: "test",
				Chat: &telebot.Chat{Type: telebot.ChatPrivate},
			},
			true,
		},
		{
			"Private chat, feeling chatty",
			1,
			"@bot",
			&telebot.Message{
				Text: "test",
				Chat: &telebot.Chat{Type: telebot.ChatPrivate},
			},
			true,
		},
		{
			"Group chat, not mentioned directly",
			0,
			"@bot",
			&telebot.Message{
				Text:     "test test test",
				Entities: []telebot.MessageEntity{{}},
				Chat:     &telebot.Chat{Type: ""},
			},
			false,
		},
		{
			"Group chat, someone else mentioned directly",
			0,
			"@bot",
			&telebot.Message{
				Text: "test test test",
				Entities: []telebot.MessageEntity{{
					Type:   telebot.EntityMention,
					Offset: 4,
					Length: 4,
				}},
				Chat: &telebot.Chat{Type: ""},
			},
			false,
		},
		{
			"Group chat, mentioned directly",
			0,
			"@bot",
			&telebot.Message{
				Text: "test @bot test",
				Entities: []telebot.MessageEntity{{
					Type:   telebot.EntityMention,
					Offset: 5,
					Length: 4,
				}},
				Chat: &telebot.Chat{Type: ""},
			},
			true,
		},
		{
			"Group chat, mentioned directly amongst others",
			0,
			"@bot",
			&telebot.Message{
				Text: "@test @bot test",
				Entities: []telebot.MessageEntity{
					{
						Type:   telebot.EntityMention,
						Offset: 0,
						Length: 5,
					},
					{
						Type:   telebot.EntityMention,
						Offset: 6,
						Length: 4,
					},
				},
				Chat: &telebot.Chat{Type: ""},
			},
			true,
		},
	}

	for _, table := range tables {
		t.Logf("Testing: %s", table.testcase)
		got := decideWhetherToRespond(table.m, table.chattiness, table.name)
		if !reflect.DeepEqual(got, table.expected) {
			t.Errorf("expected: %#v, got: %#v", table.expected, got)
		} else {
			t.Log("Passed")
		}
	}
}
