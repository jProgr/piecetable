package piecetable

import (
    "strings"
)

type PieceTable interface {
    Insert(str string, index int) PieceTable
    Append(str string) PieceTable
    DeleteChar(index int) PieceTable
    Delete(startIndex, length int) PieceTable
    ToString() string
}

// Table holds the original and text to be added.
// They are stored as pieces that point to either buffer
// (original, addendum). To get the full text, pieces are
// put together in order.
type Table struct {
    original string
    addendum strings.Builder
    pieces   []*piece
}

func NewTable(originalText string) PieceTable {
    pieces := make([]*piece, 0)

    // If originalText is not empty, create a unique piece
    // that points to all of it.
    if textLength := len(originalText); textLength != 0 {
        pieces = append(pieces, newPiece(0, textLength, originalType))
    }

    return &Table{
        originalText,
        strings.Builder{},
        pieces,
    }
}

// Append adds str to the end of the text.
func (table *Table) Append(str string) PieceTable {
    table.
        addAddendumPiece(str).
        appendToAddendum(str)

    return table
}

// Insert inserts str at index. index points to the position
// of the full text.
func (table *Table) Insert(str string, index int) PieceTable {
    // Get the piece that will be affected.
    affectedPiece := table.getPieceAtIndex(index)
    // If the index is equal or higher than the length of the
    // full text, then str is appended at the end.
    if affectedPiece == nil {
        return table.Append(str)
    }

    // Append the new text and create a piece.
    addendumPiece := table.newAddendumPiece(str)
    table.appendToAddendum(str)

    // The affectedPiece will be replaced by the contents of newPieces.
    // 3 is worst case: split the affected piece into 2 plus new piece.
    newPieces := make([]*piece, 0, 3)

    // Check whether the new piece will be added before, after, or
    // in the middle of the affected piece.

    if isBeforePiece(index, affectedPiece) {
        newPieces = append(
            newPieces,
            addendumPiece,
            affectedPiece.piece,
        )
    }

    if isAfterPiece(index, affectedPiece) {
        newPieces = append(
            newPieces,
            affectedPiece.piece,
            addendumPiece,
        )
    }

    if isMiddleOfPiece(index, affectedPiece) {
        beforePieceLength := index - affectedPiece.from
        beforePiece := newPiece(
            affectedPiece.piece.start,
            beforePieceLength,
            affectedPiece.piece.buffer,
        )
        afterPiece := newPiece(
            affectedPiece.piece.start+beforePieceLength,
            affectedPiece.piece.length-beforePieceLength,
            affectedPiece.piece.buffer,
        )

        newPieces = append(
            newPieces,
            beforePiece,
            addendumPiece,
            afterPiece,
        )
    }

    // Finally pieces are joined together, except affectedPiece.
    beforePieces := table.pieces[:affectedPiece.pieceNumber]
    afterPieces := table.pieces[affectedPiece.pieceNumber+1:]
    table.pieces = merge(beforePieces, newPieces, afterPieces)

    return table
}

func (table *Table) DeleteChar(index int) PieceTable {
    // Get the piece that will be affected.
    affectedPiece := table.getPieceAtIndex(index)
    // If the index is equal or higher than the length of the
    // full text, then there's nothing to delete.
    if affectedPiece == nil {
        return table
    }

    // If piece is already one char long, then just delete it.
    if affectedPiece.piece.length == 1 {
        beforePieces, afterPieces := splitSlice(affectedPiece.pieceNumber, table.pieces)
        table.pieces = merge(beforePieces, afterPieces[1:])

        return table
    }

    // Check where is the character being removed from: start,
    // end, or middle of the piece.

    // The char is removed from the start.
    if isBeforePiece(index, affectedPiece) {
        affectedPiece.piece.start += 1
        affectedPiece.piece.length -= 1

        return table
    }

    // The char is removed from the end.
    if isAfterPiece(index, affectedPiece) {
        affectedPiece.piece.length -= 1

        return table
    }

    // The char is removed from the middle.
    if isMiddleOfPiece(index, affectedPiece) {
        beforePiece, afterPiece := affectedPiece.piece.splitAt(index - affectedPiece.from + 1)
        beforePiece.length -= 1

        beforePieces, afterPieces := splitSlice(affectedPiece.pieceNumber, table.pieces)
        firstPieces := make([]*piece, 0, len(beforePieces)+2)
        copy(firstPieces, beforePieces)
        firstPieces = append(firstPieces, beforePiece, afterPiece)
        afterPieces = afterPieces[1:]

        table.pieces = merge(firstPieces, afterPieces)

        return table
    }

    return table
}

func (table *Table) Delete(startIndex, length int) PieceTable {
    for i := 0; i < length; i++ {
        table.DeleteChar(startIndex)
    }

    return table
}

// ToString builds the final text by applying all the pieces
// one by one.
func (table *Table) ToString() string {
    var text strings.Builder
    originalText := []rune(table.original)
    addendumText := []rune(table.addendum.String())

    for _, piece := range table.pieces {
        var str string
        start := piece.start
        end := piece.start + piece.length

        switch piece.buffer {
        case originalType:
            str = string(originalText[start:end])
        case addendumType:
            str = string(addendumText[start:end])
        }

        text.WriteString(str)
    }

    return text.String()
}

// getLengthOfAddendum returns the current length of the string in
// table.addendum.
func (table *Table) getLengthOfAddendum() int {
    return table.addendum.Len()
}

// getPieceAtIndex returns the piece that points to the full text at
// index.
//
// If index is greater than the full text, nil is returned.
func (table *Table) getPieceAtIndex(index int) *relativePiece {
    cursorPos := 0
    for i, piece := range table.pieces {
        cursorPos += piece.length
        if cursorPos > index {
            return newRelativePiece(cursorPos-piece.length, i, piece)
        }
    }

    return nil
}

// appendToAddendum appends str to table.addendum.
func (table *Table) appendToAddendum(str string) *Table {
    table.addendum.WriteString(str)

    return table
}

// addAddendumPiece creates a new piece of type addendumType and appends
// it to table.pieces.
func (table *Table) addAddendumPiece(str string) *Table {
    table.pieces = append(
        table.pieces,
        table.newAddendumPiece(str),
    )

    return table
}

// newAddendumPiece creates a new piece of type addendumType.
func (table *Table) newAddendumPiece(str string) *piece {
    return newPiece(table.getLengthOfAddendum(), len(str), addendumType)
}

// merge joins together the elements of slices into a new slice.
func merge[T any](slices ...[]T) []T {
    totalLength := 0
    for _, slice := range slices {
        totalLength += len(slice)
    }

    mergedSlice := make([]T, 0, totalLength)
    for _, slice := range slices {
        mergedSlice = append(mergedSlice, slice...)
    }

    return mergedSlice
}

func splitSlice[T any](index int, slice []T) ([]T, []T) {
    return slice[:index], slice[index:]
}

func isBeforePiece(index int, affectedPiece *relativePiece) bool {
    return index == affectedPiece.from
}

func isAfterPiece(index int, affectedPiece *relativePiece) bool {
    endIndex := affectedPiece.from + affectedPiece.piece.length

    return index == endIndex
}

func isMiddleOfPiece(index int, affectedPiece *relativePiece) bool {
    endIndex := affectedPiece.from + affectedPiece.piece.length

    return index > affectedPiece.from && index < endIndex
}
