package sheets

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"go.uber.org/zap"
)

var cellAddressRegexp = regexp.MustCompile(`^[A-Z]+[1-9]\d*$`)

func GetCellCoords(cell string) (int32, int32, error) {
	if !cellAddressRegexp.MatchString(cell) {
		logger.Errorf("Invalid cell address: %v", cell)
		return 0, 0, errors.New("invalid-cell-address-format")
	}

	rowString := ""
	columnString := ""
	for _, char := range cell {
		if char >= 'A' && char <= 'Z' {
			columnString += string(char)
		} else {
			rowString += string(char)
		}
	}

	return rowToNumber(rowString), columnToNumber(columnString), nil
}

func columnToNumber(column string) int32 {
	var result int32 = 0

	for _, char := range column {
		charNumber := int32(char) - 'A'
		result = result*26 + int32(charNumber) + 1
	}

	return result - 1
}

func rowToNumber(row string) int32 {
	result, err := strconv.ParseInt(row, 10, 32)
	if err != nil {
		logger.Error(zap.Error(err))
	}
	return int32(result - 1)
}

func GetCoordsCell(row int32, column int32) string {
	return numberToColumn(column) + numberToRow(row)
}

func numberToColumn(number int32) string {
	result := ""
	currentNumber := number + 1
	for ok := true; ok; ok = currentNumber > 0 {
		result = string(rune(currentNumber%26+'A'-1)) + result
		currentNumber = currentNumber / 26
	}
	return result
}

func numberToRow(number int32) string {
	return fmt.Sprint(number + 1)
}
