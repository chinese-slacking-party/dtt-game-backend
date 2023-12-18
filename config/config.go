package config

import (
	"time"
)

var (
	ReplicateAPIKey = "invalid"
)

var (
	CookieLife = 168 * time.Hour
)

var Outfits = []string{
	"red jacket",
	"yellow gown",
	"blue shirt",
	"brown coat",
}

var PassMsg = map[string][]string{
	"0": {
		"Amazing work! You get one point for that.",
		"Once you get 4 points, you can get a reward from your parents.",
		"Want to keep the excitement going?",
	},
}

var FailMsg = map[string][]string{
	"1": {
		"No worries at all!",
		"Every attempt is one step closer to success.",
		"Would you like to have another try?",
	},
	"2": {
		"You are doing well.",
		"I really like the effort that you are putting into this.",
		"Would you like to have another try?",
	},
	"3": {
		"No worries at all!",
		"Please try to observe their eyes, nose, mouth, and hair.",
		"These may help you find the answer.",
		"Would you like to have another try?",
	},
	"4": {
		"You are doing well.",
		"First, you can focus on the face.",
		"Second, you can focus on the eyes.",
		"Would you like to have another try?",
	},
}
