package metadata

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
)

func TestNext(t *testing.T) {
	var raw = `<?xml version="1.0" encoding="utf-8"?>
    <RETS ReplyCode="0" ReplyText="Operation successful.">
    <METADATA>
    <METADATA-CLASS Version="01.72.11582" Date="2016-03-29T21:50:11" Resource="Agent">
    </METADATA-CLASS>
    <METADATA-CLASS Version="01.72.11583" Date="2016-03-29T21:50:11" Resource="Office">
    </METADATA-CLASS>
    <METADATA-CLASS Version="01.72.11584" Date="2016-03-29T21:50:11" Resource="Listing">
    </METADATA-CLASS>
    </METADATA>
    </RETS>`
	body := ioutil.NopCloser(strings.NewReader(raw))
	defer body.Close()

	extractor := &Extractor{Body: body}
	response, err := extractor.Open()

	testutils.Ok(t, err)
	testutils.Equals(t, "Operation successful.", response.ReplyText)

	next := func(resource, version, date string) func(*testing.T) {
		return func(tt *testing.T) {
			mclass := &MClass{}
			err = extractor.DecodeNext("METADATA-CLASS", mclass)
			testutils.Ok(t, err)
			testutils.Equals(tt, resource, string(mclass.Resource))
			testutils.Equals(tt, version, string(mclass.Version))
			testutils.Equals(tt, date, string(mclass.Date))
			testutils.Equals(tt, 0, len(mclass.Class))
		}
	}

	t.Run("agent", next("Agent", "01.72.11582", "2016-03-29T21:50:11"))
	t.Run("offfice", next("Office", "01.72.11583", "2016-03-29T21:50:11"))
	t.Run("listing", next("Listing", "01.72.11584", "2016-03-29T21:50:11"))

	err = extractor.DecodeNext("METADATA-CLASS", &MClass{})
	testutils.Equals(t, io.EOF, err)
}
