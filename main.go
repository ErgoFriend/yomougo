package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
)

func main() {
	user := getUser("/960535")
	fmt.Printf("user: %+v\n", user)

	user = getUser("/1017640")
	bytes, _ := json.Marshal(&user)
	fmt.Printf("user: %s\n", string(bytes))
}

// UserInfo Cloud Functions
func UserInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		user := getUser(r.URL.Path)
		bytes, _ := json.Marshal(&user)
		fmt.Fprint(w, string(bytes))
	default:
		http.Error(w, "505 method not allowed", http.StatusMethodNotAllowed)
	}
}

// User ユーザー
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Furigana string `json:"furigana"`
	Site     string `json:"site"`
	URL      string `json:"url"`
	Content  string `json:"content"`
}

func getUser(idPath string) User {
	var baseURL = "https://mypage.syosetu.com" + idPath // likes "/1017640"

	var (
		user        User
		oneTitleAgo string
		twoTitleAgo string
	)

	geziyor.NewGeziyor(&geziyor.Options{
		LogDisabled: true,
		StartURLs:   []string{baseURL},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("dl.profile > dt, dl.profile > dd, dl.profile > dd > a").Each(func(i int, s *goquery.Selection) {

				if 0 < i {
					switch oneTitleAgo {
					case "ユーザＩＤ":
						user.ID = s.Text()
					case "ユーザネーム":
						user.Name = s.Text()
					case "フリガナ":
						user.Furigana = s.Text()
					case "サイト":
						twoTitleAgo = oneTitleAgo
						user.Site = s.Text()
						user.URL = s.AttrOr("href", "")
					case "自己紹介":
						user.Content = s.Text()
					default:
						if twoTitleAgo == "サイト" {
							user.URL = s.AttrOr("href", "")[11:]
							twoTitleAgo = ""
						}
					}
				}

				// fmt.Printf("%duser: %+v\n", i, user)
				oneTitleAgo = s.Text()
				// aaa := s.AttrOr("href", "")
				// fmt.Print(aaa[10:])

			})
		},
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
	}).Start()

	// fmt.Printf("user: %+v\n", user)
	return user
}
