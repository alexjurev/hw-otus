//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStatNegative(t *testing.T) {
	// non-valid json here - we miss } at the end
	data := `{"Id":1,"Email":"aliquid_qui_ea@Browsedrive.gov"`

	t.Run("find 'com'", func(t *testing.T) {
		_, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.Error(t, err)
	})
}

func BenchmarkGetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := `{"Id":1,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":2,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":3,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":4,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":5,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":11,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":12,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":13,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":14,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":15,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":21,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":22,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":23,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":24,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":25,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":31,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":32,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":33,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":34,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":35,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":41,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":42,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":43,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":44,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":45,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":51,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":52,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":53,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":54,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":55,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":111,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":112,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":113,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":114,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":115,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":121,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":122,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":123,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":124,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":125,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":131,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":132,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":133,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":134,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":135,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":141,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":142,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":143,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":144,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":145,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":151,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":152,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":153,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":154,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":155,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":211,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":212,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":213,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":214,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":215,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":221,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":222,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":223,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":224,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":225,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":231,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":232,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":233,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":234,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":235,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":241,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":242,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":243,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":244,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":245,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":251,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":252,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":253,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":254,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":255,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":311,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":312,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":313,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":314,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":315,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":321,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":322,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":323,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":324,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":325,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":331,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":332,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":333,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":334,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":335,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":341,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":342,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":343,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":344,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":345,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":351,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":352,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":353,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":354,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":355,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":411,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":412,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":413,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":414,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":415,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":421,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":422,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":423,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":424,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":425,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":431,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":432,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":433,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":434,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":435,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":441,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":442,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":443,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":444,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":445,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":451,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":452,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":453,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":454,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":455,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":511,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":512,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":513,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":514,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":515,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":521,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":522,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":523,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":524,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":525,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":531,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":532,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":533,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":534,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":535,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":541,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":542,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":543,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":544,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":545,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":551,"Email":"aliquid_qui_ea@Browsedrive.gov","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":552,"Email":"mLynch@broWsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":553,"Email":"RoseSmith@Browsecat.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":554,"Email":"5Moore@Teklist.net","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":555,"Email":"nulla@Linktype.com","Name":"Janice Rose","Username":"KeithHart","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`
		_, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(b, err)
	}
}
