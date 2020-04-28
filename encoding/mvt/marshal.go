package mvt

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/paulmach/orb/encoding/mvt/vectortile"
	"github.com/paulmach/orb/geojson"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

// MarshalGzipped will marshal the layers into Mapbox Vector Tile format
// and gzip the result. A lot of times MVT data is gzipped at rest,
// e.g. in a mbtiles file.
func MarshalGzipped(layers Layers) ([]byte, error) {
	data, err := Marshal(layers)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	gzwriter := gzip.NewWriter(buf)

	_, err = gzwriter.Write(data)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to write gz data")
	}

	err = gzwriter.Close()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to close gzwriter")
	}

	return buf.Bytes(), nil
}

// Marshal will take a set of layers and encode them into a Mapbox Vector Tile format.
func Marshal(layers Layers) ([]byte, error) {
	vt := &vectortile.Tile{
		Layers: make([]*vectortile.Tile_Layer, 0, len(layers)),
	}

	for _, l := range layers {
		v := l.Version
		e := l.Extent
		layer := &vectortile.Tile_Layer{
			Name:     &l.Name,
			Version:  &v,
			Extent:   &e,
			Features: make([]*vectortile.Tile_Feature, 0, len(l.Features)),
		}

		kve := newKeyValueEncoder()
		for i, f := range l.Features {
			t, g, err := encodeGeometry(f.Geometry)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf("layer %s: feature %d: error encoding geometry", l.Name, i))
			}

			tags, err := encodeProperties(kve, f.Properties)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf("layer %s: feature %d: error encoding properties", l.Name, i))
			}

			layer.Features = append(layer.Features, &vectortile.Tile_Feature{
				Id:       convertID(f.ID),
				Tags:     tags,
				Type:     &t,
				Geometry: g,
			})
		}

		layer.Keys = kve.Keys
		layer.Values = kve.Values

		vt.Layers = append(vt.Layers, layer)
	}

	return proto.Marshal(vt)
}

// UnmarshalGzipped takes gzipped Mapbox Vector Tile (MVT) data and unzips it
// before decoding it into a set of layers, It does not project the coordinates.
func UnmarshalGzipped(data []byte) (Layers, error) {
	gzreader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create gzreader")
	}

	decoded, err := ioutil.ReadAll(gzreader)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to unzip")
	}

	return Unmarshal(decoded)
}

// Unmarshal takes Mapbox Vector Tile (MVT) data and converts into a
// set of layers, It does not project the coordinates.
func Unmarshal(data []byte) (Layers, error) {
	vt := &vectortile.Tile{}
	err := vt.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	return decode(vt)
}

func decode(vt *vectortile.Tile) (Layers, error) {
	result := make(Layers, 0, len(vt.Layers))
	for i, l := range vt.Layers {
		layer := &Layer{
			Name:     l.GetName(),
			Version:  l.GetVersion(),
			Extent:   l.GetExtent(),
			Features: make([]*geojson.Feature, 0, len(l.Features)),
		}

		for j, f := range l.Features {
			geom, err := decodeGeometry(f.GetType(), f.Geometry)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf("layer %d: feature %d", i, j))
			}

			properties := decodeFeatureProperties(l.Keys, l.Values, f.Tags)

			if geom != nil {
				gjf := &geojson.Feature{
					Geometry:   geom,
					Properties: properties,
				}

				if f.Id != nil {
					gjf.ID = float64(*f.Id)
				}

				layer.Features = append(layer.Features, gjf)
			}
		}

		result = append(result, layer)
	}

	return result, nil
}

func encodeProperties(kve *keyValueEncoder, properties geojson.Properties) ([]uint32, error) {
	tags := make([]uint32, 0, 2*len(properties))
	for k, v := range properties {
		ki := kve.Key(k)
		vi, err := kve.Value(v)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("property %s", k))
		}

		tags = append(tags, ki, vi)
	}

	return tags, nil
}

func decodeFeatureProperties(
	keys []string,
	values []*vectortile.Tile_Value,
	tags []uint32,
) geojson.Properties {
	result := make(geojson.Properties, len(tags)/2)
	if len(tags) == 0 {
		return result
	}

	for i := 2; i <= len(tags); i += 2 {
		vi := tags[i-1]
		if int(vi) >= len(values) {
			continue
		}

		v := decodeValue(values[vi])
		if v != nil {
			ti := tags[i-2]
			if int(ti) >= len(keys) {
				continue
			}

			result[keys[ti]] = v
		}
	}

	return result
}

func convertID(id interface{}) *uint64 {
	if id == nil {
		return nil
	}

	switch id := id.(type) {
	case int:
		return convertIntID(id)
	case int8:
		return convertIntID(int(id))
	case int16:
		return convertIntID(int(id))
	case int32:
		return convertIntID(int(id))
	case int64:
		return convertIntID(int(id))
	case uint:
		v := uint64(id)
		return &v
	case uint8:
		v := uint64(id)
		return &v
	case uint16:
		v := uint64(id)
		return &v
	case uint32:
		v := uint64(id)
		return &v
	case uint64:
		v := uint64(id)
		return &v
	case float32:
		return convertIntID(int(id))
	case float64:
		return convertIntID(int(id))
	case string:
		i, err := strconv.Atoi(id)
		if err == nil {
			return convertIntID(i)
		}
	}

	return nil
}

func convertIntID(i int) *uint64 {
	if i < 0 {
		return nil
	}

	v := uint64(i)
	return &v
}
