package engine

import (
	"github.com/mumax/3/cuda"
	"github.com/mumax/3/data"
	"github.com/mumax/3/script"
	"github.com/mumax/3/util"
	"reflect"
)

// Special buffered quantity to store magnetization
// makes sure it's normalized etc.
type magnetization struct {
	buffer_ *data.Slice
}

func (m *magnetization) Mesh() *data.Mesh     { return Mesh() }
func (m *magnetization) NComp() int           { return 3 }
func (m *magnetization) Name() string         { return "m" }
func (m *magnetization) Unit() string         { return "" }
func (m *magnetization) Buffer() *data.Slice  { return m.buffer_ }
func (m *magnetization) Child() []script.Expr { return nil }

// allocate storage (not done by init, as mesh size may not yet be known then)
func (m *magnetization) alloc() {
	m.buffer_ = cuda.NewSlice(3, m.Mesh().Size())
	m.Set(RandomMag()) // sane starting config
}

func (b *magnetization) SetArray(src *data.Slice) {
	if src.Size() != b.Mesh().Size() {
		src = data.Resample(src, b.Mesh().Size())
	}
	data.Copy(b.Buffer(), src)
	cuda.Normalize(b.Buffer(), geometry.Gpu())
}

func (m *magnetization) Set(c Config) {
	m.SetInShape(nil, c)
}

func (m *magnetization) LoadFile(fname string) {
	m.SetArray(LoadFile(fname))
}

func (m *magnetization) Slice() (s *data.Slice, recycle bool) {
	return m.Buffer(), false
}

func (m *magnetization) Region(r int) *sliceInRegion { return &sliceInRegion{m, r} }

func (m *magnetization) String() string { return util.Sprint(m.Buffer().HostCopy()) }

// Set the value of one cell.
func (m *magnetization) SetCell(ix, iy, iz int, v ...float64) {
	nComp := m.NComp()
	util.Argument(len(v) == nComp)
	for c := 0; c < nComp; c++ {
		cuda.SetCell(m.Buffer(), c, ix, iy, iz, float32(v[c]))
	}
}

// Get the value of one cell.
func (m *magnetization) GetCell(comp, ix, iy, iz int) float64 {
	return float64(cuda.GetCell(m.Buffer(), comp, ix, iy, iz))
}

func (m *magnetization) TableData() []float64 { return Average(m) }

// Sets the magnetization inside the shape
func (m *magnetization) SetInShape(region Shape, conf Config) {
	if region == nil {
		region = universe
	}
	host := m.Buffer().HostCopy()
	h := host.Vectors()
	n := m.Mesh().Size()

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				r := Index2Coord(ix, iy, iz)
				x, y, z := r[X], r[Y], r[Z]
				if region(x, y, z) { // inside
					m := conf(x, y, z)
					h[X][iz][iy][ix] = float32(m[X])
					h[Y][iz][iy][ix] = float32(m[Y])
					h[Z][iz][iy][ix] = float32(m[Z])
				}
			}
		}
	}
	m.SetArray(host)
}

// set m to config in region
func (m *magnetization) SetRegion(region int, conf Config) {
	host := m.Buffer().HostCopy()
	h := host.Vectors()
	n := m.Mesh().Size()
	r := byte(region)

	for iz := 0; iz < n[Z]; iz++ {
		for iy := 0; iy < n[Y]; iy++ {
			for ix := 0; ix < n[X]; ix++ {
				pos := Index2Coord(ix, iy, iz)
				x, y, z := pos[X], pos[Y], pos[Z]
				if regions.arr[iz][iy][ix] == r {
					m := conf(x, y, z)
					h[X][iz][iy][ix] = float32(m[X])
					h[Y][iz][iy][ix] = float32(m[Y])
					h[Z][iz][iy][ix] = float32(m[Z])
				}
			}
		}
	}
	m.SetArray(host)
}

func (m *magnetization) SetValue(v interface{})  { m.SetInShape(nil, v.(Config)) }
func (m *magnetization) InputType() reflect.Type { return reflect.TypeOf(Config(nil)) }
func (m *magnetization) Type() reflect.Type      { return reflect.TypeOf(new(magnetization)) }
func (m *magnetization) Eval() interface{}       { return m }
func (m *magnetization) Average() data.Vector    { return unslice(Average(&M)) }

func normalize(m *data.Slice) {
	cuda.Normalize(m, geometry.Gpu())
}

func (m *magnetization) resize(s2 [3]int) {
	backup := m.Buffer().HostCopy()
	resized := data.Resample(backup, s2)
	m.buffer_.Free()
	m.buffer_ = cuda.NewSlice(VECTOR, s2)
	data.Copy(m.buffer_, resized)
}
