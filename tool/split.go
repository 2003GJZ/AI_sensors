package tool

import "strings"

func SplitString(ags string, partition string) (string, string, error) { //返回截取字符串（去除分割符），和截取完的字符串
	index := strings.Index(ags, partition)
	if index == -1 {
		return ags, "", nil
	}
	return ags[:index], ags[index+len(partition):], nil
}
