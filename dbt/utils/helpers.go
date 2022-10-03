package utils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func Contains(stringList []string, value string) bool {
	for _, val := range stringList {
		if val == value {
			return true
		}
	}
	return false
}

func RemoveFromList(originalList []string, itemsToRemove []string) []string {
	var returnList []string
	returnList = append(returnList, originalList...)
	for _, item := range itemsToRemove {
		returnList = findAndDeleteFirstItem(returnList, item)
	}

	return returnList
}

func findAndDeleteFirstItem(stringList []string, item string) []string {
	for index, i := range stringList {
		if i == item {
			return append(stringList[:index], stringList[index+1:]...)
		}
	}

	var dst []string
	return append(dst, stringList...)
}

func InterfaceToStringList(rawInterface interface{}) []string {
	rawList := rawInterface.(*schema.Set).List()
	stringList := make([]string, len(rawList))
	for i, v := range rawList {
		stringList[i] = v.(string)
	}
	return stringList
}
