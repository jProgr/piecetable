package piecetable

type bufferType int

const (
    originalType bufferType = iota
    addendumType
)

type piece struct {
    start  int
    length int
    buffer bufferType
}

func newPiece(start, length int, buffer bufferType) *piece {
    return &piece{start, length, buffer}
}

func (piece piece) splitAt(length int) (*piece, *piece) {
    bufferType := piece.buffer
    beforePiece := newPiece(piece.start, length, bufferType)
    afterPiece := newPiece(piece.start+length, piece.length-length, bufferType)

    return beforePiece, afterPiece
}

type relativePiece struct {
    from        int
    pieceNumber int
    piece       *piece
}

func newRelativePiece(from, pieceNumber int, piece *piece) *relativePiece {
    return &relativePiece{from, pieceNumber, piece}
}
