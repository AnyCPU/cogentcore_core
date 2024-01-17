// Code generated by "core generate"; DO NOT EDIT.

package fi

import (
	"errors"
	"log"
	"strconv"

	"cogentcore.org/core/enums"
)

var _CatValues = []Cat{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

// CatN is the highest valid value
// for type Cat, plus one.
const CatN Cat = 16

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the enumgen command to generate them again.
func _CatNoOp() {
	var x [1]struct{}
	_ = x[UnknownCat-(0)]
	_ = x[Folder-(1)]
	_ = x[Archive-(2)]
	_ = x[Backup-(3)]
	_ = x[Code-(4)]
	_ = x[Doc-(5)]
	_ = x[Sheet-(6)]
	_ = x[Data-(7)]
	_ = x[Text-(8)]
	_ = x[Image-(9)]
	_ = x[Model-(10)]
	_ = x[Audio-(11)]
	_ = x[Video-(12)]
	_ = x[Font-(13)]
	_ = x[Exe-(14)]
	_ = x[Bin-(15)]
}

var _CatNameToValueMap = map[string]Cat{
	`UnknownCat`: 0,
	`Folder`:     1,
	`Archive`:    2,
	`Backup`:     3,
	`Code`:       4,
	`Doc`:        5,
	`Sheet`:      6,
	`Data`:       7,
	`Text`:       8,
	`Image`:      9,
	`Model`:      10,
	`Audio`:      11,
	`Video`:      12,
	`Font`:       13,
	`Exe`:        14,
	`Bin`:        15,
}

var _CatDescMap = map[Cat]string{
	0:  `UnknownCat is an unknown file category`,
	1:  `Folder is a folder / directory`,
	2:  `Archive is a collection of files, e.g., zip tar`,
	3:  `Backup is a backup file (# ~ etc)`,
	4:  `Code is a programming language file`,
	5:  `Doc is an editable word processing file including latex, markdown, html, css, etc`,
	6:  `Sheet is a spreadsheet file (.xls etc)`,
	7:  `Data is some kind of data format (csv, json, database, etc)`,
	8:  `Text is some other kind of text file`,
	9:  `Image is an image (jpeg, png, svg, etc) *including* PDF`,
	10: `Model is a 3D model`,
	11: `Audio is an audio file`,
	12: `Video is a video file`,
	13: `Font is a font file`,
	14: `Exe is a binary executable file (scripts go in Code)`,
	15: `Bin is some other type of binary (object files, libraries, etc)`,
}

var _CatMap = map[Cat]string{
	0:  `UnknownCat`,
	1:  `Folder`,
	2:  `Archive`,
	3:  `Backup`,
	4:  `Code`,
	5:  `Doc`,
	6:  `Sheet`,
	7:  `Data`,
	8:  `Text`,
	9:  `Image`,
	10: `Model`,
	11: `Audio`,
	12: `Video`,
	13: `Font`,
	14: `Exe`,
	15: `Bin`,
}

// String returns the string representation
// of this Cat value.
func (i Cat) String() string {
	if str, ok := _CatMap[i]; ok {
		return str
	}
	return strconv.FormatInt(int64(i), 10)
}

// SetString sets the Cat value from its
// string representation, and returns an
// error if the string is invalid.
func (i *Cat) SetString(s string) error {
	if val, ok := _CatNameToValueMap[s]; ok {
		*i = val
		return nil
	}
	return errors.New(s + " is not a valid value for type Cat")
}

// Int64 returns the Cat value as an int64.
func (i Cat) Int64() int64 {
	return int64(i)
}

// SetInt64 sets the Cat value from an int64.
func (i *Cat) SetInt64(in int64) {
	*i = Cat(in)
}

// Desc returns the description of the Cat value.
func (i Cat) Desc() string {
	if str, ok := _CatDescMap[i]; ok {
		return str
	}
	return i.String()
}

// CatValues returns all possible values
// for the type Cat.
func CatValues() []Cat {
	return _CatValues
}

// Values returns all possible values
// for the type Cat.
func (i Cat) Values() []enums.Enum {
	res := make([]enums.Enum, len(_CatValues))
	for i, d := range _CatValues {
		res[i] = d
	}
	return res
}

// IsValid returns whether the value is a
// valid option for type Cat.
func (i Cat) IsValid() bool {
	_, ok := _CatMap[i]
	return ok
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i Cat) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *Cat) UnmarshalText(text []byte) error {
	if err := i.SetString(string(text)); err != nil {
		log.Println("Cat.UnmarshalText:", err)
	}
	return nil
}

var _KnownValues = []Known{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124}

// KnownN is the highest valid value
// for type Known, plus one.
const KnownN Known = 125

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the enumgen command to generate them again.
func _KnownNoOp() {
	var x [1]struct{}
	_ = x[Unknown-(0)]
	_ = x[Any-(1)]
	_ = x[AnyKnown-(2)]
	_ = x[AnyFolder-(3)]
	_ = x[AnyArchive-(4)]
	_ = x[Multipart-(5)]
	_ = x[Tar-(6)]
	_ = x[Zip-(7)]
	_ = x[GZip-(8)]
	_ = x[SevenZ-(9)]
	_ = x[Xz-(10)]
	_ = x[BZip-(11)]
	_ = x[Dmg-(12)]
	_ = x[Shar-(13)]
	_ = x[AnyBackup-(14)]
	_ = x[Trash-(15)]
	_ = x[AnyCode-(16)]
	_ = x[Ada-(17)]
	_ = x[Bash-(18)]
	_ = x[Csh-(19)]
	_ = x[C-(20)]
	_ = x[CSharp-(21)]
	_ = x[D-(22)]
	_ = x[Diff-(23)]
	_ = x[Eiffel-(24)]
	_ = x[Erlang-(25)]
	_ = x[Forth-(26)]
	_ = x[Fortran-(27)]
	_ = x[FSharp-(28)]
	_ = x[Go-(29)]
	_ = x[Haskell-(30)]
	_ = x[Java-(31)]
	_ = x[JavaScript-(32)]
	_ = x[Lisp-(33)]
	_ = x[Lua-(34)]
	_ = x[Makefile-(35)]
	_ = x[Mathematica-(36)]
	_ = x[Matlab-(37)]
	_ = x[ObjC-(38)]
	_ = x[OCaml-(39)]
	_ = x[Pascal-(40)]
	_ = x[Perl-(41)]
	_ = x[Php-(42)]
	_ = x[Prolog-(43)]
	_ = x[Python-(44)]
	_ = x[R-(45)]
	_ = x[Ruby-(46)]
	_ = x[Rust-(47)]
	_ = x[Scala-(48)]
	_ = x[Tcl-(49)]
	_ = x[AnyDoc-(50)]
	_ = x[BibTeX-(51)]
	_ = x[TeX-(52)]
	_ = x[Texinfo-(53)]
	_ = x[Troff-(54)]
	_ = x[Html-(55)]
	_ = x[Css-(56)]
	_ = x[Markdown-(57)]
	_ = x[Rtf-(58)]
	_ = x[MSWord-(59)]
	_ = x[OpenText-(60)]
	_ = x[OpenPres-(61)]
	_ = x[MSPowerpoint-(62)]
	_ = x[EBook-(63)]
	_ = x[EPub-(64)]
	_ = x[AnySheet-(65)]
	_ = x[MSExcel-(66)]
	_ = x[OpenSheet-(67)]
	_ = x[AnyData-(68)]
	_ = x[Csv-(69)]
	_ = x[Json-(70)]
	_ = x[Xml-(71)]
	_ = x[Protobuf-(72)]
	_ = x[Ini-(73)]
	_ = x[Tsv-(74)]
	_ = x[Uri-(75)]
	_ = x[Color-(76)]
	_ = x[GoGi-(77)]
	_ = x[Yaml-(78)]
	_ = x[AnyText-(79)]
	_ = x[PlainText-(80)]
	_ = x[ICal-(81)]
	_ = x[VCal-(82)]
	_ = x[VCard-(83)]
	_ = x[AnyImage-(84)]
	_ = x[Pdf-(85)]
	_ = x[Postscript-(86)]
	_ = x[Gimp-(87)]
	_ = x[GraphVis-(88)]
	_ = x[Gif-(89)]
	_ = x[Jpeg-(90)]
	_ = x[Png-(91)]
	_ = x[Svg-(92)]
	_ = x[Tiff-(93)]
	_ = x[Pnm-(94)]
	_ = x[Pbm-(95)]
	_ = x[Pgm-(96)]
	_ = x[Ppm-(97)]
	_ = x[Xbm-(98)]
	_ = x[Xpm-(99)]
	_ = x[Bmp-(100)]
	_ = x[Heic-(101)]
	_ = x[Heif-(102)]
	_ = x[AnyModel-(103)]
	_ = x[Vrml-(104)]
	_ = x[X3d-(105)]
	_ = x[AnyAudio-(106)]
	_ = x[Aac-(107)]
	_ = x[Flac-(108)]
	_ = x[Mp3-(109)]
	_ = x[Ogg-(110)]
	_ = x[Midi-(111)]
	_ = x[Wav-(112)]
	_ = x[AnyVideo-(113)]
	_ = x[Mpeg-(114)]
	_ = x[Mp4-(115)]
	_ = x[Mov-(116)]
	_ = x[Ogv-(117)]
	_ = x[Wmv-(118)]
	_ = x[Avi-(119)]
	_ = x[AnyFont-(120)]
	_ = x[TrueType-(121)]
	_ = x[WebOpenFont-(122)]
	_ = x[AnyExe-(123)]
	_ = x[AnyBin-(124)]
}

var _KnownNameToValueMap = map[string]Known{
	`Unknown`:      0,
	`Any`:          1,
	`AnyKnown`:     2,
	`AnyFolder`:    3,
	`AnyArchive`:   4,
	`Multipart`:    5,
	`Tar`:          6,
	`Zip`:          7,
	`GZip`:         8,
	`SevenZ`:       9,
	`Xz`:           10,
	`BZip`:         11,
	`Dmg`:          12,
	`Shar`:         13,
	`AnyBackup`:    14,
	`Trash`:        15,
	`AnyCode`:      16,
	`Ada`:          17,
	`Bash`:         18,
	`Csh`:          19,
	`C`:            20,
	`CSharp`:       21,
	`D`:            22,
	`Diff`:         23,
	`Eiffel`:       24,
	`Erlang`:       25,
	`Forth`:        26,
	`Fortran`:      27,
	`FSharp`:       28,
	`Go`:           29,
	`Haskell`:      30,
	`Java`:         31,
	`JavaScript`:   32,
	`Lisp`:         33,
	`Lua`:          34,
	`Makefile`:     35,
	`Mathematica`:  36,
	`Matlab`:       37,
	`ObjC`:         38,
	`OCaml`:        39,
	`Pascal`:       40,
	`Perl`:         41,
	`Php`:          42,
	`Prolog`:       43,
	`Python`:       44,
	`R`:            45,
	`Ruby`:         46,
	`Rust`:         47,
	`Scala`:        48,
	`Tcl`:          49,
	`AnyDoc`:       50,
	`BibTeX`:       51,
	`TeX`:          52,
	`Texinfo`:      53,
	`Troff`:        54,
	`Html`:         55,
	`Css`:          56,
	`Markdown`:     57,
	`Rtf`:          58,
	`MSWord`:       59,
	`OpenText`:     60,
	`OpenPres`:     61,
	`MSPowerpoint`: 62,
	`EBook`:        63,
	`EPub`:         64,
	`AnySheet`:     65,
	`MSExcel`:      66,
	`OpenSheet`:    67,
	`AnyData`:      68,
	`Csv`:          69,
	`Json`:         70,
	`Xml`:          71,
	`Protobuf`:     72,
	`Ini`:          73,
	`Tsv`:          74,
	`Uri`:          75,
	`Color`:        76,
	`GoGi`:         77,
	`Yaml`:         78,
	`AnyText`:      79,
	`PlainText`:    80,
	`ICal`:         81,
	`VCal`:         82,
	`VCard`:        83,
	`AnyImage`:     84,
	`Pdf`:          85,
	`Postscript`:   86,
	`Gimp`:         87,
	`GraphVis`:     88,
	`Gif`:          89,
	`Jpeg`:         90,
	`Png`:          91,
	`Svg`:          92,
	`Tiff`:         93,
	`Pnm`:          94,
	`Pbm`:          95,
	`Pgm`:          96,
	`Ppm`:          97,
	`Xbm`:          98,
	`Xpm`:          99,
	`Bmp`:          100,
	`Heic`:         101,
	`Heif`:         102,
	`AnyModel`:     103,
	`Vrml`:         104,
	`X3d`:          105,
	`AnyAudio`:     106,
	`Aac`:          107,
	`Flac`:         108,
	`Mp3`:          109,
	`Ogg`:          110,
	`Midi`:         111,
	`Wav`:          112,
	`AnyVideo`:     113,
	`Mpeg`:         114,
	`Mp4`:          115,
	`Mov`:          116,
	`Ogv`:          117,
	`Wmv`:          118,
	`Avi`:          119,
	`AnyFont`:      120,
	`TrueType`:     121,
	`WebOpenFont`:  122,
	`AnyExe`:       123,
	`AnyBin`:       124,
}

var _KnownDescMap = map[Known]string{
	0:   `Unknown = a non-known file type`,
	1:   `Any is used when selecting a file type, if any type is OK (including Unknown) see also AnyKnown and the Any options for each category`,
	2:   `AnyKnown is used when selecting a file type, if any Known file type is OK (excludes Unknown) -- see Any and Any options for each category`,
	3:   `Folder is a folder / directory`,
	4:   `Archive is a collection of files, e.g., zip tar`,
	5:   ``,
	6:   ``,
	7:   ``,
	8:   ``,
	9:   ``,
	10:  ``,
	11:  ``,
	12:  ``,
	13:  ``,
	14:  `Backup files`,
	15:  ``,
	16:  `Code is a programming language file`,
	17:  ``,
	18:  ``,
	19:  ``,
	20:  ``,
	21:  ``,
	22:  ``,
	23:  ``,
	24:  ``,
	25:  ``,
	26:  ``,
	27:  ``,
	28:  ``,
	29:  ``,
	30:  ``,
	31:  ``,
	32:  ``,
	33:  ``,
	34:  ``,
	35:  ``,
	36:  ``,
	37:  ``,
	38:  ``,
	39:  ``,
	40:  ``,
	41:  ``,
	42:  ``,
	43:  ``,
	44:  ``,
	45:  ``,
	46:  ``,
	47:  ``,
	48:  ``,
	49:  ``,
	50:  `Doc is an editable word processing file including latex, markdown, html, css, etc`,
	51:  ``,
	52:  ``,
	53:  ``,
	54:  ``,
	55:  ``,
	56:  ``,
	57:  ``,
	58:  ``,
	59:  ``,
	60:  ``,
	61:  ``,
	62:  ``,
	63:  ``,
	64:  ``,
	65:  `Sheet is a spreadsheet file (.xls etc)`,
	66:  ``,
	67:  ``,
	68:  `Data is some kind of data format (csv, json, database, etc)`,
	69:  ``,
	70:  ``,
	71:  ``,
	72:  ``,
	73:  ``,
	74:  ``,
	75:  ``,
	76:  ``,
	77:  ``,
	78:  ``,
	79:  `Text is some other kind of text file`,
	80:  ``,
	81:  ``,
	82:  ``,
	83:  ``,
	84:  `Image is an image (jpeg, png, svg, etc) *including* PDF`,
	85:  ``,
	86:  ``,
	87:  ``,
	88:  ``,
	89:  ``,
	90:  ``,
	91:  ``,
	92:  ``,
	93:  ``,
	94:  ``,
	95:  ``,
	96:  ``,
	97:  ``,
	98:  ``,
	99:  ``,
	100: ``,
	101: ``,
	102: ``,
	103: `Model is a 3D model`,
	104: ``,
	105: ``,
	106: `Audio is an audio file`,
	107: ``,
	108: ``,
	109: ``,
	110: ``,
	111: ``,
	112: ``,
	113: `Video is a video file`,
	114: ``,
	115: ``,
	116: ``,
	117: ``,
	118: ``,
	119: ``,
	120: `Font is a font file`,
	121: ``,
	122: ``,
	123: `Exe is a binary executable file`,
	124: `Bin is some other unrecognized binary type`,
}

var _KnownMap = map[Known]string{
	0:   `Unknown`,
	1:   `Any`,
	2:   `AnyKnown`,
	3:   `AnyFolder`,
	4:   `AnyArchive`,
	5:   `Multipart`,
	6:   `Tar`,
	7:   `Zip`,
	8:   `GZip`,
	9:   `SevenZ`,
	10:  `Xz`,
	11:  `BZip`,
	12:  `Dmg`,
	13:  `Shar`,
	14:  `AnyBackup`,
	15:  `Trash`,
	16:  `AnyCode`,
	17:  `Ada`,
	18:  `Bash`,
	19:  `Csh`,
	20:  `C`,
	21:  `CSharp`,
	22:  `D`,
	23:  `Diff`,
	24:  `Eiffel`,
	25:  `Erlang`,
	26:  `Forth`,
	27:  `Fortran`,
	28:  `FSharp`,
	29:  `Go`,
	30:  `Haskell`,
	31:  `Java`,
	32:  `JavaScript`,
	33:  `Lisp`,
	34:  `Lua`,
	35:  `Makefile`,
	36:  `Mathematica`,
	37:  `Matlab`,
	38:  `ObjC`,
	39:  `OCaml`,
	40:  `Pascal`,
	41:  `Perl`,
	42:  `Php`,
	43:  `Prolog`,
	44:  `Python`,
	45:  `R`,
	46:  `Ruby`,
	47:  `Rust`,
	48:  `Scala`,
	49:  `Tcl`,
	50:  `AnyDoc`,
	51:  `BibTeX`,
	52:  `TeX`,
	53:  `Texinfo`,
	54:  `Troff`,
	55:  `Html`,
	56:  `Css`,
	57:  `Markdown`,
	58:  `Rtf`,
	59:  `MSWord`,
	60:  `OpenText`,
	61:  `OpenPres`,
	62:  `MSPowerpoint`,
	63:  `EBook`,
	64:  `EPub`,
	65:  `AnySheet`,
	66:  `MSExcel`,
	67:  `OpenSheet`,
	68:  `AnyData`,
	69:  `Csv`,
	70:  `Json`,
	71:  `Xml`,
	72:  `Protobuf`,
	73:  `Ini`,
	74:  `Tsv`,
	75:  `Uri`,
	76:  `Color`,
	77:  `GoGi`,
	78:  `Yaml`,
	79:  `AnyText`,
	80:  `PlainText`,
	81:  `ICal`,
	82:  `VCal`,
	83:  `VCard`,
	84:  `AnyImage`,
	85:  `Pdf`,
	86:  `Postscript`,
	87:  `Gimp`,
	88:  `GraphVis`,
	89:  `Gif`,
	90:  `Jpeg`,
	91:  `Png`,
	92:  `Svg`,
	93:  `Tiff`,
	94:  `Pnm`,
	95:  `Pbm`,
	96:  `Pgm`,
	97:  `Ppm`,
	98:  `Xbm`,
	99:  `Xpm`,
	100: `Bmp`,
	101: `Heic`,
	102: `Heif`,
	103: `AnyModel`,
	104: `Vrml`,
	105: `X3d`,
	106: `AnyAudio`,
	107: `Aac`,
	108: `Flac`,
	109: `Mp3`,
	110: `Ogg`,
	111: `Midi`,
	112: `Wav`,
	113: `AnyVideo`,
	114: `Mpeg`,
	115: `Mp4`,
	116: `Mov`,
	117: `Ogv`,
	118: `Wmv`,
	119: `Avi`,
	120: `AnyFont`,
	121: `TrueType`,
	122: `WebOpenFont`,
	123: `AnyExe`,
	124: `AnyBin`,
}

// String returns the string representation
// of this Known value.
func (i Known) String() string {
	if str, ok := _KnownMap[i]; ok {
		return str
	}
	return strconv.FormatInt(int64(i), 10)
}

// SetString sets the Known value from its
// string representation, and returns an
// error if the string is invalid.
func (i *Known) SetString(s string) error {
	if val, ok := _KnownNameToValueMap[s]; ok {
		*i = val
		return nil
	}
	return errors.New(s + " is not a valid value for type Known")
}

// Int64 returns the Known value as an int64.
func (i Known) Int64() int64 {
	return int64(i)
}

// SetInt64 sets the Known value from an int64.
func (i *Known) SetInt64(in int64) {
	*i = Known(in)
}

// Desc returns the description of the Known value.
func (i Known) Desc() string {
	if str, ok := _KnownDescMap[i]; ok {
		return str
	}
	return i.String()
}

// KnownValues returns all possible values
// for the type Known.
func KnownValues() []Known {
	return _KnownValues
}

// Values returns all possible values
// for the type Known.
func (i Known) Values() []enums.Enum {
	res := make([]enums.Enum, len(_KnownValues))
	for i, d := range _KnownValues {
		res[i] = d
	}
	return res
}

// IsValid returns whether the value is a
// valid option for type Known.
func (i Known) IsValid() bool {
	_, ok := _KnownMap[i]
	return ok
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i Known) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *Known) UnmarshalText(text []byte) error {
	if err := i.SetString(string(text)); err != nil {
		log.Println("Known.UnmarshalText:", err)
	}
	return nil
}
