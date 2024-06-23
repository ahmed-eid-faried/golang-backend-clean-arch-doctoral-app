package main

// import (
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"os"
// 	"sort"
//

// 	"github.com/gorilla/pat"
// 	"github.com/markbates/goth"
// 	"github.com/markbates/goth/gothic"
// 	"github.com/markbates/goth/providers/apple"
// 	"github.com/markbates/goth/providers/facebook"
// 	"github.com/markbates/goth/providers/github"
// 	"github.com/markbates/goth/providers/google"
// 	"github.com/markbates/goth/providers/linkedin"
// 	"github.com/markbates/goth/providers/microsoftonline"
// )

// // ProviderIndex stores information about the available providers.
// type ProviderIndex struct {
// 	Providers    []string
// 	ProvidersMap map[string]string
// }

// var (
// 	// Define templates
// 	indexTemplate = `{{range $key,$value:=.Providers}}
//     <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
// {{end}}`

// 	userTemplate = `
// <p><a href="/logout/{{.Provider}}">logout</a></p>
// <p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
// <p>Email: {{.Email}}</p>
// <p>NickName: {{.NickName}}</p>
// <p>Location: {{.Location}}</p>
// <p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
// <p>Description: {{.Description}}</p>
// <p>UserID: {{.UserID}}</p>
// <p>AccessToken: {{.AccessToken}}</p>
// <p>ExpiresAt: {{.ExpiresAt}}</p>
// <p>RefreshToken: {{.RefreshToken}}</p>
// `
// )

// // SetupProviders initializes the authentication providers.
// func SetupProviders() {
// 	goth.UseProviders(
// 		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
// 		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
// 		apple.New(os.Getenv("APPLE_KEY"), os.Getenv("APPLE_SECRET"), "http://localhost:3000/auth/apple/callback", nil, apple.ScopeName, apple.ScopeEmail),
// 		linkedin.New(os.Getenv("LINKEDIN_KEY"), os.Getenv("LINKEDIN_SECRET"), "http://localhost:3000/auth/linkedin/callback"),
// 		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
// 		microsoftonline.New(os.Getenv("MICROSOFTONLINE_KEY"), os.Getenv("MICROSOFTONLINE_SECRET"), "http://localhost:3000/auth/microsoftonline/callback"),
// 	)
// }

// // CreateProviderIndex creates a sorted list of providers and their display names.
// func CreateProviderIndex() *ProviderIndex {
// 	m := map[string]string{
// 		"google":          "Google",
// 		"facebook":        "Facebook",
// 		"apple":           "Apple",
// 		"linkedin":        "LinkedIn",
// 		"github":          "GitHub",
// 		"microsoftonline": "Microsoft",
// 	}
// 	var keys []string
// 	for k := range m {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)

// 	return &ProviderIndex{Providers: keys, ProvidersMap: m}
// }

// // Handler sets up the routes and starts the HTTP server.
// func Handler() {
// 	providerIndex := CreateProviderIndex()

// 	p := pat.New()
// 	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
// 		user, err := gothic.CompleteUserAuth(res, req)
// 		if err != nil {
// 			fmt.Fprintln(res, err)
// 			return
// 		}
// 		t, _ := template.New("userTemplate").Parse(userTemplate)
// 		t.Execute(res, user)
// 	})

// 	p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
// 		gothic.Logout(res, req)
// 		res.Header().Set("Location", "/")
// 		res.WriteHeader(http.StatusTemporaryRedirect)
// 	})

// 	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
// 		if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
// 			t, _ := template.New("userTemplate").Parse(userTemplate)
// 			t.Execute(res, gothUser)
// 		} else {
// 			gothic.BeginAuthHandler(res, req)
// 		}
// 	})

// 	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
// 		t, _ := template.New("indexTemplate").Parse(indexTemplate)
// 		t.Execute(res, providerIndex)
// 	})

// 	log.Println("listening on localhost:3000")
// 	log.Fatal(http.ListenAndServe(":3000", p))
// }
