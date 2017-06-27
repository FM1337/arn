package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixUserNick(t *testing.T) {
	// Nickname autocorrect
	assert.True(t, FixUserNick("Akyoto") == "Akyoto")
	assert.True(t, FixUserNick("Tsundere") == "Tsundere")
	assert.True(t, FixUserNick("akyoto") == "Akyoto")
	assert.True(t, FixUserNick("aky123oto") == "Akyoto")
	assert.True(t, FixUserNick("__aky123oto%$§") == "Akyoto")
	assert.True(t, FixUserNick("__aky123oto%$§__") == "Akyoto")
	assert.True(t, FixUserNick("123%&/(__%") == "")
}

func TestFixAccountNick(t *testing.T) {
	// Nickname autocorrect
	assert.True(t, FixAccountNick("UserName") == "UserName")
	assert.True(t, FixAccountNick("anilist.co/user/UserName") == "UserName")
	assert.True(t, FixAccountNick("https://anilist.co/user/UserName") == "UserName")
	assert.True(t, FixAccountNick("osu.ppy.sh/u/UserName") == "UserName")
	assert.True(t, FixAccountNick("kitsu.io/users/UserName/library") == "UserName")
}

func TestFixTag(t *testing.T) {
	// Nickname autocorrect
	assert.Equal(t, FixTag("general"), "general")
	assert.Equal(t, FixTag("https://notify.moe/anime/244"), "anime:244")
	assert.Equal(t, FixTag("https://notify.moe/anime/244/"), "anime:244")
	assert.Equal(t, FixTag("https://osu.ppy.sh/s/320118"), "osu-beatmap:320118")
}
