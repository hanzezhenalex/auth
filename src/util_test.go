package src

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func _sort(s []string) []string {
	SortSliceAsc(s)
	return s
}

func TestSliceAppend(t *testing.T) {
	rq := require.New(t)

	t.Run("append", func(t *testing.T) {
		s1 := []string{"1", "3", "5", "7"}
		s2 := []string{"0", "2", "4", "6", "8"}

		s3, duplicated := SliceAppend(s1, s2)
		rq.EqualValues([]string{
			"0", "1", "2", "3", "4", "5", "6", "7", "8",
		}, _sort(s3))
		rq.Equal(0, len(duplicated))
	})

	t.Run("duplicated", func(t *testing.T) {
		s1 := []string{"1", "3", "5", "7", "8"}
		s2 := []string{"0", "2", "4", "6", "8"}

		s3, duplicated := SliceAppend(s1, s2)
		rq.EqualValues([]string{
			"0", "1", "2", "3", "4", "5", "6", "7", "8",
		}, _sort(s3))
		rq.EqualValues([]string{"8"}, duplicated)
	})
}

func TestSliceRemove(t *testing.T) {
	rq := require.New(t)

	t.Run("remove", func(t *testing.T) {
		s1 := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8"}
		s2 := []string{"0", "2", "4", "6", "8"}

		s3, nonExisted := SliceRemove(s1, s2)
		rq.EqualValues([]string{"1", "3", "5", "7"}, _sort(s3))
		rq.Equal(0, len(nonExisted))
	})

	t.Run("non-existed", func(t *testing.T) {
		s1 := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8"}
		s2 := []string{"0", "2", "4", "6", "9"}

		s3, nonExisted := SliceRemove(s1, s2)
		rq.EqualValues([]string{"1", "3", "5", "7", "8"}, _sort(s3))
		rq.EqualValues([]string{"9"}, nonExisted)
	})
}
