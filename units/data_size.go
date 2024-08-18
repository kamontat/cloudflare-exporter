package units

type DataSize int64

func (d DataSize) Byte() int {
	return int(d)
}
func (d DataSize) KiloByte() int {
	return int(d / KiloByte)
}
func (d DataSize) Megabyte() int {
	return int(d / Megabyte)
}
func (d DataSize) Gigabyte() int {
	return int(d / Gigabyte)
}

func (d DataSize) String() string {
	var arr [32]byte
	format(int64(d), &arr, unitMap, []string{"GB", "MB", "KB", "B"})
	return string(arr[:])
}

const (
	Byte     DataSize = 1
	KiloByte DataSize = 1024 * Byte
	Megabyte DataSize = 1024 * KiloByte
	Gigabyte DataSize = 1024 * Megabyte

	KiloByteSI DataSize = 1000 * Byte
	MegabyteSI DataSize = 1000 * KiloByteSI
	GigabyteSI DataSize = 1000 * MegabyteSI
)

var unitMap = map[string]uint64{
	"B":  uint64(Byte),
	"KB": uint64(KiloByte),
	"MB": uint64(Megabyte),
	"GB": uint64(Gigabyte),
}

func ParseDataSize(in string) (DataSize, error) {
	return parseUnit(in, unitMap, func(number int64) DataSize {
		return DataSize(number)
	})
}
