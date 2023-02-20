package os

import (
	"os/user"
	"strings"
)

func GetUserName() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}

	if strings.Index(u.Username, "\\") < 0 {
		return u.Username
	}

	us := strings.Split(u.Username, "\\")

	if len(us) > 1 {
		return us[len(us) - 1]
	}

	return u.Username
}
