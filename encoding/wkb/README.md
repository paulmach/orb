encoding/wkb [![Godoc Reference](https://godoc.org/github.com/paulmach/orb?status.svg)](https://godoc.org/github.com/paulmach/orb/encoding/wkb)
============

This package provides encoding and decoding of [WKB](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry#Well-known_binary)
data. The interface is defined as:

	func Marshal(geom orb.Geometry, byteOrder ...binary.ByteOrder) ([]byte, error)
	func MustMarshal(geom orb.Geometry, byteOrder ...binary.ByteOrder) []byte

	func NewEncoder(w io.Writer) *Encoder
	func (e *Encoder) SetByteOrder(bo binary.ByteOrder)
	func (e *Encoder) Encode(geom orb.Geometry) error

	func Unmarshal(b []byte) (orb.Geometry, error)

	func NewDecoder(r io.Reader) *Decoder
	func (d *Decoder) Decode() (orb.Geometry, error)

### Reading and Writing to a SQL database

This package provides wrappers for `orb.Geometry` types that implement
`sql.Scanner` and `driver.Value`. For example:

	row := db.QueryRow("SELECT ST_AsBinary(point_column) FROM postgis_table")

	var p orb.Point
	err := row.Scan(wkb.Scanner(&p))

	db.Exec("INSERT INTO table (point_column) VALUES (?)",
		wkb.Value(p))

If you don't know the type of the geometry try something like

	s := wkb.Scanner(nil)
	err := row.Scan(&s)

	switch g := s.Geometry.(type) {
	case orb.Point:
	case orb.LineString:
	}

Scanning directly from MySQL columns is supported. By default MySQL returns geometry
data as WKB but prefixed with a 4 byte SRID. To support this, if the data is not
valid WKB, the code will strip the first 4 bytes, the SRID, and try again.
This works for most use cases.
