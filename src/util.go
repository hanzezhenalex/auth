package src

import (
	"crypto/rand"
	"encoding/base64"
	"sort"
)

// GenerateSecureRandomString 生成指定长度的安全随机字符串
func GenerateSecureRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 使用 base64 进行编码，确保生成的字符串是可打印的
	randomString := base64.StdEncoding.EncodeToString(randomBytes)

	// 去除可能的末尾填充字符
	return randomString[:length], nil
}

func SortSliceAsc(slices ...[]string) {
	for _, s := range slices {
		sort.Slice(s, func(i, j int) bool {
			return s[i] < s[j]
		})
	}
}

func sliceOnDiffItems(
	master []string,
	participate []string,
	onSameItem func(int, int),
	onMasterUnique func(int),
	onParticipateUnique func(int),
) {
	SortSliceAsc(master, participate)

	i, j := -1, 0

	for ; j < len(participate); j++ {
		for i = i + 1; i < len(master) && master[i] < participate[j]; i++ {
			if onMasterUnique != nil {
				onMasterUnique(i)
			}
		}

		if i >= len(master) {
			onParticipateUnique(j)
			continue
		} else if master[i] == participate[j] {
			if onSameItem != nil {
				onSameItem(i, j)
			}
		} else {
			if onParticipateUnique != nil {
				onParticipateUnique(j)
			}
		}
	}

	if onMasterUnique != nil {
		for i = i + 1; i < len(master); i++ {
			onMasterUnique(i)
		}
	}
}

func SliceAppend(origin []string, add []string) ([]string, []string) {
	var duplicated []string
	sliceOnDiffItems(origin, add,
		func(_ int, j int) {
			duplicated = append(duplicated, add[j])
		},
		nil,
		func(j int) {
			origin = append(origin, add[j])
		},
	)
	return origin, duplicated
}

func SliceRemove(origin []string, remove []string) ([]string, []string) {
	var removed, nonExisted []string
	sliceOnDiffItems(origin, remove,
		nil,
		func(i int) {
			removed = append(removed, origin[i])
		},
		func(j int) {
			nonExisted = append(nonExisted, remove[j])
		},
	)
	return removed, nonExisted
}
