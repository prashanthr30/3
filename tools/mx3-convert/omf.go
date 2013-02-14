package main

// OVF2 suport
// Author: Mykola Dvornik
// Modified by Arne Vansteenkiste, 2011, 2012, 2013.

import (
	"code.google.com/p/mx3/data"
	"fmt"
	"io"
	"log"
	"strings"
	"unsafe"
)

func dumpOmf(out io.Writer, file string, q *data.Slice, dataformat string) {

	switch strings.ToLower(dataformat) {
	case "binary", "binary 4":
		dataformat = "Binary 4"
	case "text":
		dataformat = "Text"
	default:
		log.Fatalf("Illegal OMF data format: %v", dataformat)
	}

	writeOmfHeader(out, q)
	writeOmfData(out, q, dataformat)
	hdr(out, "End", "Segment")
}

const (
	OMF_CONTROL_NUMBER = 1234567.0 // The omf format requires the first encoded number in the binary data section to be this control number
)

func writeOmfData(out io.Writer, q *data.Slice, dataformat string) {

	hdr(out, "Begin", "Data "+dataformat)
	switch strings.ToLower(dataformat) {
	case "text":
		writeOmfText(out, q)
	case "binary 4":
		writeOmfBinary4(out, q)
	default:
		log.Fatalf("Illegal OMF data format: %v. Options are: Text, Binary 4", dataformat)
	}
	hdr(out, "End", "Data "+dataformat)
}

// Writes the OMF header
func writeOmfHeader(out io.Writer, q *data.Slice) {
	gridsize := q.Mesh().Size()
	cellsize := q.Mesh().CellSize()

	hdr(out, "OOMMF", "rectangular mesh v1.0")
	hdr(out, "Segment count", "1")
	hdr(out, "Begin", "Segment")

	hdr(out, "Begin", "Header")

	dsc(out, "Time", 0) //q.Time) // TODO !!
	hdr(out, "Title", q.Tag())
	hdr(out, "meshtype", "rectangular")
	hdr(out, "meshunit", "m")
	hdr(out, "xbase", cellsize[Z]/2)
	hdr(out, "ybase", cellsize[Y]/2)
	hdr(out, "zbase", cellsize[X]/2)
	hdr(out, "xstepsize", cellsize[Z])
	hdr(out, "ystepsize", cellsize[Y])
	hdr(out, "zstepsize", cellsize[X])
	hdr(out, "xmin", 0)
	hdr(out, "ymin", 0)
	hdr(out, "zmin", 0)
	hdr(out, "xmax", cellsize[Z]*float64(gridsize[Z]))
	hdr(out, "ymax", cellsize[Y]*float64(gridsize[Y]))
	hdr(out, "zmax", cellsize[X]*float64(gridsize[X]))
	hdr(out, "xnodes", gridsize[Z])
	hdr(out, "ynodes", gridsize[Y])
	hdr(out, "znodes", gridsize[X])
	hdr(out, "ValueRangeMinMag", 1e-08) // not so "optional" as the OOMMF manual suggests...
	hdr(out, "ValueRangeMaxMag", 1)     // TODO
	hdr(out, "valueunit", "?")
	hdr(out, "valuemultiplier", 1)

	hdr(out, "End", "Header")
}

// Writes data in OMF Binary 4 format
func writeOmfBinary4(out io.Writer, array *data.Slice) {
	data := array.Tensors()
	gridsize := array.Mesh().Size()

	var bytes []byte

	// OOMMF requires this number to be first to check the format
	var controlnumber float32 = OMF_CONTROL_NUMBER
	// Conversion form float32 [4]byte in big-endian
	// Inlined for performance, terabytes of data will pass here...
	bytes = (*[4]byte)(unsafe.Pointer(&controlnumber))[:]
	bytes[0], bytes[1], bytes[2], bytes[3] = bytes[3], bytes[2], bytes[1], bytes[0] // swap endianess
	out.Write(bytes)

	// Here we loop over X,Y,Z, not Z,Y,X, because
	// internal in C-order == external in Fortran-order
	ncomp := array.NComp()
	for i := 0; i < gridsize[X]; i++ {
		for j := 0; j < gridsize[Y]; j++ {
			for k := 0; k < gridsize[Z]; k++ {
				for c := 0; c < ncomp; c++ {
					// dirty conversion from float32 to [4]byte
					bytes = (*[4]byte)(unsafe.Pointer(&data[SwapIndex(c, ncomp)][i][j][k]))[:]
					bytes[0], bytes[1], bytes[2], bytes[3] = bytes[3], bytes[2], bytes[1], bytes[0]
					out.Write(bytes)
				}
			}
		}
	}
}

// Writes data in OMF Text format
func writeOmfText(out io.Writer, tens *data.Slice) (err error) {

	data := tens.Tensors()
	gridsize := tens.Mesh().Size()

	// Here we loop over X,Y,Z, not Z,Y,X, because
	// internal in C-order == external in Fortran-order
	for i := 0; i < gridsize[X]; i++ {
		for j := 0; j < gridsize[Y]; j++ {
			for k := 0; k < gridsize[Z]; k++ {
				for c := 0; c < tens.NComp(); c++ {
					_, err = fmt.Fprint(out, data[SwapIndex(c, tens.NComp())][i][j][k], " ") // converts to user space.
				}
				_, err = fmt.Fprint(out, "\n")
			}
		}
	}
	return
}

//func floats2bytes(floats []float32) []byte {
//	return (*[4]byte)(unsafe.Pointer(&floats[0]))[:]
//}

// Writes a header key/value pair to out:
// # Key: Value
func hdr(out io.Writer, key string, value ...interface{}) {
	fmt.Fprint(out, "# ", key, ": ")
	fmt.Fprintln(out, value...)
}

func dsc(out io.Writer, k, v interface{}) {
	hdr(out, "Desc", k, ": ", v)
}
