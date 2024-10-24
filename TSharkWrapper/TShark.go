package TSharkWrapper

import (
	"log"
	"os/exec"
)

type TShark struct {
	tsharkCmd *exec.Cmd
}

func NewTShark() *TShark {
	cmd := exec.Command("tshark", "-e ip.addr -e frame.len -i any -T fields")
	return &TShark{tsharkCmd: cmd}
}

func (tshark *TShark) Run() {
	err := tshark.tsharkCmd.Run()
	if err != nil {
		log.Fatal("Unable to start tshark due to ", err)
	}
}

func (tshark *TShark) RedirectOutputToChannelAsync(unParsedDataChannel chan string) {
	for {
		bytes, err := tshark.tsharkCmd.Output()
		if err != nil {
			log.Fatal("Unable to get the output of tshark due to: ", err.Error())
		}

		if len(unParsedDataChannel) < cap(unParsedDataChannel) {
			unParsedDataChannel <- string(bytes)
		}
	}
}
