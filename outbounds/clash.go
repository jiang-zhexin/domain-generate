package outbounds

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"Domain_Generate/data"
	"Domain_Generate/log"
)

func save2clash(dtl *data.DomainList, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Warn("无法创建目录：%s", dir)
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		log.Warn("无法创建文件：%s", path)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for domain := range dtl.Suffix {
		_, err := writer.WriteString("DOMAIN-SUFFIX," + domain + "\n")
		if err != nil {
			return err
		}
	}
	for domain := range dtl.Full {
		_, err := writer.WriteString("DOMAIN," + domain + "\n")
		if err != nil {
			return err
		}
	}
	for domain := range dtl.Keyword {
		_, err := writer.WriteString("DOMAIN-KEYWORD," + domain + "\n")
		if err != nil {
			return err
		}
	}
	// clash 不支持正则域名

	err = writer.Flush()
	if err != nil {
		log.Warn(fmt.Sprintf("写入文件：%s 错误", path))
		return err
	}
	return nil
}
