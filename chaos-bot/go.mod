module ultre.me/radio-chaos/chaos-bot

go 1.14

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/etherlabsio/errors v0.2.3
	github.com/etherlabsio/pkg v0.0.0-20191020161600-58998d98f9ce
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/httplog v0.2.0
	github.com/gohugoio/hugo v0.83.1
	github.com/graarh/golang-socketio v0.0.0-20170510162725-2c44953b9b5f
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/oklog/run v1.1.0
	github.com/peterbourgon/ff/v3 v3.0.0
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/tpyolang/tpyo-cli v1.0.0
	github.com/ultreme/histoire-pour-enfant-generator v0.0.0-20200402084311-66b2cd0d2da6
	gopkg.in/yaml.v2 v2.4.0
	moul.io/godev v1.7.0
	moul.io/moulsay v1.3.0
	moul.io/number-to-words v0.6.0
	moul.io/pipotron v1.13.4
	ultre.me/recettator v0.4.0
	ultre.me/smsify v1.0.0
)

// replace github.com/ultreme/histoire-pour-enfant-generator => ../../histoire-pour-enfant-generator
