package gfp

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/tlianza/surf"
	"github.com/tlianza/surf/browser"
)

type Exercise struct {
	FitocracyName string `toml:"fitocracy_name"`
	FitocracyId   int    `toml:"fitocracy_id"`
	MFPName       string `toml:"mfp_name"`
	MFPId         int    `toml:"mfp_id"`
}

type Config struct {
	Test      string
	Exercises []Exercise `toml:"exercises"`
}

type Fetcher struct {
	host    string
	browser *browser.Browser
}

func NewFetcher() *Fetcher {
	var config Config
	_, err := toml.DecodeFile("/go/src/github.com/tlianza/gfp/gfp.toml", &config)
	if nil != err {
		fmt.Println(err)
	}

	fmt.Println("loaded:")
	fmt.Println(config)

	return &Fetcher{
		host:    "www.myfitnesspal.com",
		browser: surf.NewBrowser(),
	}
}

func (f *Fetcher) Login(username, password string) error {
	err := f.browser.Open(fmt.Sprintf("https://%s/account/login", f.host))
	if err != nil {
		return err
	}

	// Log in to the site.
	fm, _ := f.browser.Form("form.LoginForm")
	fm.Input("username", username)
	fm.Input("password", password)
	return fm.Submit()
}

func (f *Fetcher) AddExercise(exercise_id, date, sets, reps, weight string) error {
	err := f.browser.Open(fmt.Sprintf("http://%s/exercise/search?type=%s", f.host, "strength"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Opened Exercise search: " + f.browser.Title())

	fm, err := f.browser.Form("form#add_exercise")
	if err != nil {
		panic(err)
	}

	token := f.browser.Dom().Find("input[name='authenticity_token']").First()
	tokenValue, _ := token.Attr("value")
	fm.Set("utf8", "✓")
	fm.Set("authenticity_token", tokenValue)
	fm.Set("calorie_multiplier", "1.0")
	fm.Set("search", "bench press")
	fm.Set("exercise_entry[exercise_id]", exercise_id)
	fm.Set("exercise_entry[date]", date)
	fm.Set("exercise_entry[exercise_type]", "1")
	fm.Set("authenticity_token", tokenValue)
	fm.Set("exercise_entry[sets]", sets)
	fm.Set("exercise_entry[quantity]", reps)
	fm.Set("exercise_entry[weight]", weight)

	return fm.Submit()
}
