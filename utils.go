package terminal

import (
)


func min(a, b int) int{
    if a > b {
        return b
    }
    return a
}

func max(a, b int) int{
    if a > b {
        return a
    }
    return b
}

// func ArrayRemove(slice []interface{},  i int) []interface{} {
    // return append(slice[:i], slice[i + 1:]...)
// }

// func Remove(slice []int, s int) []int {
    // return append(slice[:s], slice[s+1:]...)
// }

// func Insert(slice []int, i int, item int) []int{
    // newSlice := make([]int, 0)
    // if i > 0 {
        // newSlice = append(newSlice, slice[:i]...)
    // }

    // newSlice = append(newSlice, item)
    // fmt.Println(newSlice)
    // newSlice = append(newSlice, slice[i:]...)
    // fmt.Println(newSlice)
    // return newSlice
// }
