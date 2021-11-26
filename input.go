package toold

import (
	"bufio"
	"fmt"
	"os"
)

/*
SysScanInputText 系统扫描等待输入
*/
func SysScanInputText(prompt string, sysFunc func(text string) string) {
	go func(prompt string) {
		input := bufio.NewScanner(os.Stdin)
		promptCopy := prompt
		for {
			if len(promptCopy) != 0 {
				fmt.Printf("%v", promptCopy)
			}
			if input.Scan() {
				test := input.Text()
				promptText := sysFunc(test)
				if len(promptText) != 0 {
					promptCopy = promptText
				}
			}
		}
	}(prompt)
}
