package bot

import (
	"fmt"
	"os"
)

func (s *Svc) CreateFile(report, userID string) {
	f, err := os.Create(fmt.Sprintf("%s_report.txt", userID))
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(report)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
