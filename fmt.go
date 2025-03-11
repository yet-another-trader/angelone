package angelone

import (
	"bytes"
	"io"

	"github.com/apache/arrow/go/v18/arrow"
	"github.com/apache/arrow/go/v18/arrow/array"
	"github.com/apache/arrow/go/v18/arrow/memory"
	"github.com/apache/arrow/go/v18/parquet"
	"github.com/apache/arrow/go/v18/parquet/compress"
	"github.com/apache/arrow/go/v18/parquet/pqarrow"
)

func (c Candles) Parquet() (io.Reader, error) {
	pool := memory.NewGoAllocator()

	schema := arrow.NewSchema([]arrow.Field{
		{Name: "Time", Type: arrow.PrimitiveTypes.Int64},
		{Name: "Open", Type: arrow.PrimitiveTypes.Float64},
		{Name: "High", Type: arrow.PrimitiveTypes.Float64},
		{Name: "Low", Type: arrow.PrimitiveTypes.Float64},
		{Name: "Close", Type: arrow.PrimitiveTypes.Float64},
		{Name: "Volume", Type: arrow.PrimitiveTypes.Uint32},
	}, nil)

	bldr := array.NewRecordBuilder(pool, schema)
	defer bldr.Release()

	timeBldr := bldr.Field(0).(*array.Int64Builder)
	openBldr := bldr.Field(1).(*array.Float64Builder)
	highBldr := bldr.Field(2).(*array.Float64Builder)
	lowBldr := bldr.Field(3).(*array.Float64Builder)
	closeBldr := bldr.Field(4).(*array.Float64Builder)
	volumeBldr := bldr.Field(5).(*array.Uint32Builder)

	for _, candle := range c {
		timeBldr.Append(candle.Time.Unix())
		openBldr.Append(candle.Open.InexactFloat64())
		highBldr.Append(candle.High.InexactFloat64())
		lowBldr.Append(candle.Low.InexactFloat64())
		closeBldr.Append(candle.Close.InexactFloat64())
		volumeBldr.Append(uint32(candle.Volume))
	}

	record := bldr.NewRecord()
	defer record.Release()

	props := parquet.NewWriterProperties(
		parquet.WithCompression(compress.Codecs.Snappy),
		parquet.WithRootRepetition(parquet.Repetitions.Required),
	)

	var buf bytes.Buffer

	writer, err := pqarrow.NewFileWriter(schema, &buf, props, pqarrow.DefaultWriterProps())
	if err != nil {
		return nil, err
	}
	defer writer.Close()

	err = writer.Write(record)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
