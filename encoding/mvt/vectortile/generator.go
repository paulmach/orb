package vectortile

//go:generate protoc  --proto_path=../vectortile --go_opt=paths=source_relative --go_opt=Mvector_tile.proto=./vectortile --go_out=.  vector_tile.proto
