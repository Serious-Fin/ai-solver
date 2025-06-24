package main
import "testing"
func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num

		if prevIndex, ok := numMap[complement]; ok {
			return []int{prevIndex, i}
		}

		numMap[num] = i
	}

	return nil
}
func TestTwoSum0(t *testing.T) {
	got := twoSum([]int{2, 7, 11, 15}, 9)
	want := []int{0, 1}
	if !areEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
func TestTwoSum1(t *testing.T) {
	got := twoSum([]int{1, 2, 3, 4, 4}, 8)
	want := []int{3, 4}
	if !areEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
func areEqual(got []int, want []int) bool {
	if len(got) != 2 {
		return false
	}

	if got[0] == want[0] && got[1] == want[1] {
		return true
	}
	if got[0] == want[1] && got[1] == want[0] {
		return true
	}

	return false
}
