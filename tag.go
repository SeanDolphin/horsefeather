package horsefeather

import "strings"

type tag struct {
	Tag string
	CRC bool
}

func parseTag(rawTag string) tag {
	t := tag{
		Tag: rawTag,
		CRC: false,
	}
	if strings.Contains(rawTag, ",") {
		parts := strings.Split(rawTag, ",")
		if len(parts) > 1 {
			t.Tag = parts[0]
			for _, part := range parts {
				switch part {
				case "crc":
					t.CRC = true
				}
			}
		}
	}
	return t
}
