package autocorrect_test

import (
	"testing"

	"github.com/animenotifier/arn/autocorrect"
	"github.com/stretchr/testify/assert"
)

func TestFixUserNick(t *testing.T) {
	// Nickname autocorrect
	assert.True(t, autocorrect.UserNick("Akyoto") == "Akyoto")
	assert.True(t, autocorrect.UserNick("Tsundere") == "Tsundere")
	assert.True(t, autocorrect.UserNick("akyoto") == "Akyoto")
	assert.True(t, autocorrect.UserNick("aky123oto") == "Akyoto")
	assert.True(t, autocorrect.UserNick("__aky123oto%$§") == "Akyoto")
	assert.True(t, autocorrect.UserNick("__aky123oto%$§__") == "Akyoto")
	assert.True(t, autocorrect.UserNick("123%&/(__%") == "")
}

func TestFixAccountNick(t *testing.T) {
	// Nickname autocorrect
	assert.True(t, autocorrect.AccountNick("UserName") == "UserName")
	assert.True(t, autocorrect.AccountNick("anilist.co/user/UserName") == "UserName")
	assert.True(t, autocorrect.AccountNick("https://anilist.co/user/UserName") == "UserName")
	assert.True(t, autocorrect.AccountNick("osu.ppy.sh/u/UserName") == "UserName")
	assert.True(t, autocorrect.AccountNick("kitsu.io/users/UserName/library") == "UserName")
}

func TestFixTag(t *testing.T) {
	// Nickname autocorrect
	assert.Equal(t, autocorrect.Tag("general"), "general")
	assert.Equal(t, autocorrect.Tag("https://notify.moe/anime/244"), "anime:244")
	assert.Equal(t, autocorrect.Tag("https://notify.moe/anime/244/"), "anime:244")
	assert.Equal(t, autocorrect.Tag("https://osu.ppy.sh/s/320118"), "osu-beatmap:320118")
}
