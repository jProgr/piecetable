package piecetable

import (
    "fmt"
    "testing"
)

func TestAddsTextToTheEndOnEmptyOriginalText(test *testing.T) {
    table := NewTable("")
    table.Append("abc")

    result := table.ToString()
    if result != "abc" {
        test.Fatal("Expected: abc. Actual: " + fmt.Sprintf("%v", result))
    }
}

func TestCreatesATableWithNoEmptyOriginalText(test *testing.T) {
    text := "abc"
    result := NewTable(text).ToString()
    if result != text {
        test.Fatal(fmt.Sprintf("Expected: %v. Actual: %v.", text, result))
    }
}

func TestAddsTextAtTheStart(test *testing.T) {
    table := NewTable("abc")
    table.Insert("xyz ", 0)

    result := table.ToString()
    if expected := "xyz abc"; result != expected {
        test.Fatal(fmt.Sprintf("Expected: %v. Actual: %v.", expected, result))
    }
}

func TestAddsTextAtTheEnd(test *testing.T) {
    table := NewTable("abc")
    table.Insert(" xyz", 3)

    result := table.ToString()
    if expected := "abc xyz"; result != expected {
        test.Fatal(fmt.Sprintf("Expected: %v. Actual: %v.", expected, result))
    }
}

func TestAddsTextInTheMiddle(test *testing.T) {
    table := NewTable("abc def")
    table.Insert("xyz ", 4)

    result := table.ToString()
    if expected := "abc xyz def"; result != expected {
        test.Fatal(fmt.Sprintf("Expected: %v. Actual: %v.", expected, result))
    }
}

func TestAddsTextMultipleTimes(test *testing.T) {
    table := NewTable("")
    table.
        Insert("abc", 0).
        Insert(" ghi", 3).
        Insert(" def", 3)

    result := table.ToString()
    if expected := "abc def ghi"; result != expected {
        test.Fatal(fmt.Sprintf("Expected: %v. Actual: %v.", expected, result))
    }
}

func TestDoesNothingWhenTryingToDeleteCharOutsideOfText(test *testing.T) {
    text := "abc"
    table := NewTable(text)

    assertEquals(text, table.DeleteChar(100).ToString(), test)
}

func TestDeletesOneCharAtTheStart(test *testing.T) {
    table := NewTable("abc")

    assertEquals("bc", table.DeleteChar(0).ToString(), test)
}

func TestDeletesOneCharAtTheEnd(test *testing.T) {
    table := NewTable("abc")

    assertEquals("ab", table.DeleteChar(2).ToString(), test)
}

func TestDeletesOneCharInTheMiddle(test *testing.T) {
    table := NewTable("abc")

    assertEquals("ac", table.DeleteChar(1).ToString(), test)
}

func TestDeletesOneCharInItsOwnPiece(test *testing.T) {
    table := NewTable("abc")

    assertEquals(
        "a",
        table.
            DeleteChar(1).
            DeleteChar(1).
            ToString(),
        test,
    )
}

func TestDeletesMultipleChars(test *testing.T) {
    table := NewTable("abc").Insert(" xyz", 3)

    assertEquals("az", table.Delete(1, 5).ToString(), test)
}

func assertEquals[U comparable](expected, actual U, test *testing.T) {
    if expected != actual {
        test.Fatal(fmt.Sprintf("Expected: %v. Actual: %v.", expected, actual))
    }
}
