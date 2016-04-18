package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	// "os"
	// "sort"
)

type ReplacementEntry struct {
	Original string
	New_str  string
}

func get_distance(str1, str2 string) int {
	return DistanceForStrings([]rune(str1), []rune(str2), DefaultOptions)
}

func distance_matrix(arr1, arr2 []string) [][]int {
	distance_mat := make([][]int, len(arr1))

	for i := 0; i < len(arr1); i++ {
		distance_mat[i] = make([]int, len(arr2))
	}

	for idx1, str1 := range arr1 {
		for idx2, str2 := range arr2 {
			distance_mat[idx1][idx2] = get_distance(str1, str2)
		}
	}

	return distance_mat
}

// Gets the distance for each array then organizes the returned array into
// an array of strings that is in the order
// TODO: MAke this not garbage
func order_by_distance(arr1, arr2 []string) ([]string, []string) {
	var smallerArr, biggerArr []string

	if len(arr1) <= len(arr2) {
		smallerArr = arr1
		biggerArr = arr2
	} else {
		smallerArr = arr2
		biggerArr = arr1
	}

	// only ever be as long as the smallest array
	var orderedArr1 = make([]string, len(smallerArr))
	var orderedArr2 = make([]string, len(smallerArr))

	distance_mat := distance_matrix(smallerArr, biggerArr)

	fmt.Println(distance_mat)
	var closest, closestIdx int

	for idx1, str1 := range smallerArr {
		closest = distance_mat[idx1][0]
		closestIdx = 0
		for idx2, _ := range biggerArr {
			if distance_mat[idx1][idx2] <= closest {
				closest = distance_mat[idx1][idx2]
				closestIdx = idx2
			}
		}

		orderedArr1[idx1] = str1
		orderedArr2[idx1] = biggerArr[closestIdx]
	}

	return orderedArr1, orderedArr2
}

func normalize_comparision(arr1, arr2 []string) []string {
	orderedArr1, orderedArr2 := order_by_distance(arr1, arr2)
	fmt.Printf("Arr1 = %+v\n Arr2 = %+v", orderedArr1, orderedArr2)
	normlized_arr := make([]string, 0, len(orderedArr1))

	var tmp_arr []string
	// if len(arr1) != len(arr2) {
	// 	fmt.Printf("Not supported yet")
	// 	return normlized_arr
	// }

	for idx, str1 := range orderedArr1 {
		str2 := orderedArr2[idx]
		if str1 == str2 {
			normlized_arr = append(normlized_arr, str1)
			continue
		}
		tmp_arr = FindAllSubstrings(str1, str2, 5)
		if len(tmp_arr) != 1 {
			fmt.Printf("Something is wrong %+v", tmp_arr)
			continue
		}

		normlized_arr = append(normlized_arr, tmp_arr[0])
	}

	return normlized_arr
}

func SampleSubstrings(dirty_strings []string) []string {
	var sub_strs, comp_sub_strs []string

	// common_sub_strings := make([]string, 0, len(dirty_strings))

	comp_sub_strs = FindAllSubstrings(dirty_strings[0], dirty_strings[1], 10)

	fmt.Printf("%+v", comp_sub_strs)

	idx_offset := 2
	for idx, _ := range dirty_strings[idx_offset:] {
		fmt.Printf("INdexing %d\n", idx)
		sub_strs = FindAllSubstrings(dirty_strings[idx+idx_offset-1], dirty_strings[idx+idx_offset], 10)
		fmt.Printf("Using %+v | %+v\n", dirty_strings[idx+idx_offset-1], dirty_strings[idx+idx_offset])
		fmt.Printf("Strs\n")
		fmt.Printf("%+v\n", comp_sub_strs)
		fmt.Printf("%+v\n", sub_strs)
		fmt.Printf("======\n")

		comp_sub_strs = normalize_comparision(sub_strs, comp_sub_strs)
	}

	return comp_sub_strs
}

// Returns all common sequential substrings
// Based on
// https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Longest_common_substring#Go
func FindAllSubstrings(s1 string, s2 string, threshold int) []string {
	var similarity_matrix = make([][]int, 1+len(s1))
	var sub_strings = map[string]string{}

	var key string
	var prev_key string

	for i := 0; i < len(similarity_matrix); i++ {
		similarity_matrix[i] = make([]int, 1+len(s2))
	}

	for x := 1; x < 1+len(s1); x++ {
		for y := 1; y < 1+len(s2); y++ {
			if s1[x-1] == s2[y-1] {
				similarity_matrix[x][y] = similarity_matrix[x-1][y-1] + 1
				if similarity_matrix[x][y] >= threshold {
					key = strconv.Itoa(x) + "," + strconv.Itoa(y)
					prev_key = strconv.Itoa(x-1) + "," + strconv.Itoa(y-1)

					delete(sub_strings, prev_key)
					sub_strings[key] = s1[x-similarity_matrix[x][y] : x]
				}
			} else {
				similarity_matrix[x][y] = 0
			}
		}
	}

	sub_strings_as_arr := string_map_to_arr(sub_strings)
	return sub_strings_as_arr
}

func string_map_to_arr(string_map map[string]string) []string {
	as_arr := make([]string, 0, len(string_map))

	for _, value := range string_map {
		as_arr = append(as_arr, value)
	}

	return as_arr
}

// Removes any string from the given map
func RemoveSubStrings(str_arr []string, common_sub_strs []string, threshold int) []string {
	var new_strs = make([]string, len(str_arr))
	var new_str string

	for idx, str := range str_arr {
		new_str = str

		for _, common_sub_str := range common_sub_strs {
			new_str = strings.Replace(new_str, common_sub_str, "", -1)
		}

		new_strs[idx] = new_str
	}

	return new_strs
}

func RemoveStringMatch(strings []string, match_to_remove *regexp.Regexp) []string {
	var new_strs = make([]string, len(strings))
	var start, end int
	var str_index []int
	var new_str string

	for idx, str := range strings {
		str_index = match_to_remove.FindStringIndex(str)
		fmt.Printf("\nRemoving %s match %+v\n", str, str_index)
		if len(str_index) < 2 {
			continue
		}
		start, end = str_index[0], str_index[1]

		new_str = str[:start] + str[end:] // slice out the match
		new_strs[idx] = new_str
	}

	return new_strs
}

func getDirContents(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	var file_names = make([]string, len(files))
	if err != nil {
		log.Fatal(err)
	}

	for idx, file := range files {
		file_names[idx] = file.Name()
	}

	return file_names
}

func TestLCS() []ReplacementEntry {
	// episodes := []string{
	// 	"Community.S02E01.Anthropology.101.720p.WEB-DL.x264-LeRalouf.mkv",
	// 	"Community.S02E02.Accounting.For.Lawyers.720p.WEB-DL.x264-LeRalouf.mkv",
	// }

	episodes := getDirContents("test")
	fmt.Printf("%v\n", episodes)
	var replacements = make([]ReplacementEntry, len(episodes))

	re := regexp.MustCompile("(S?\\d{1,2})(E?\\d{2})")
	cleaned_strs := RemoveStringMatch(episodes, re)

	// common_sub_strs := FindAllSubstrings(cleaned_strs[0], cleaned_strs[1], 4)
	common_sub_strs := SampleSubstrings(cleaned_strs)

	// FindAllSubstrings(cleaned_strs[0], cleaned_strs[1], 4)

	new_strings := RemoveSubStrings(cleaned_strs, common_sub_strs, 3)
	for idx, episode := range episodes {
		replacements[idx].Original = episode
		replacements[idx].New_str = new_strings[idx]
	}

	fmt.Printf("%+v\n", replacements)
	replace_map_str, _ := json.Marshal(&replacements)
	fmt.Printf("%s\n", replace_map_str)
	return replacements
}
