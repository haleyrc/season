package season

import (
	"os"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestFiles(t *testing.T) {
	testcases := map[string]struct {
		base string
		opts []Option
		want FileMods
	}{
		"flat": {
			base: "./testdata/flat",
			opts: nil,
			want: FileMods{
				base:    "./testdata/flat",
				seasons: map[string]int{},
				mods: []*FileMod{
					&FileMod{
						base:   "./testdata/flat",
						subdir: "",
						from:   "1. First Video.mp4",
						toName: "S01E01_First_Video",
						toExt:  ".mp4",
					},
					&FileMod{
						base:   "./testdata/flat",
						subdir: "",
						from:   "2. Second Video.mp4",
						toName: "S01E02_Second_Video",
						toExt:  ".mp4",
					},
				},
			},
		},
		"nested": {
			base: "./testdata/nested",
			opts: []Option{
				WithNested(true),
			},
			want: FileMods{
				base: "./testdata/nested",
				seasons: map[string]int{
					"1. First Season":  1,
					"2. Second Season": 1,
				},
				mods: []*FileMod{
					&FileMod{
						base:   "./testdata/nested",
						subdir: "1. First Season",
						from:   "1. First Video.mp4",
						toName: "S01E01_First_Video",
						toExt:  ".mp4",
					},
					&FileMod{
						base:   "./testdata/nested",
						subdir: "2. Second Season",
						from:   "1. First Video.mp4",
						toName: "S02E01_First_Video",
						toExt:  ".mp4",
					},
				},
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			got, err := ScanV2(tc.base, tc.opts...)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("=== ScanV2(%s)\n\ngot:\n%s\n\nwant:\n%s", tc.base, spew.Sdump(got), spew.Sdump(tc.want))
			}
			got.Display(os.Stdout)
		})
	}
}

// func TestDisplay(t *testing.T) {
// 	files := Files{
// 		BasePath: "./some/dir",
// 		Files: []File{
// 			File{Original: "a", Modified: "b"},
// 			File{Original: "aa", Modified: "bb"},
// 			File{Original: "aaa", Modified: "bbb"},
// 			File{Original: "aaaa", Modified: "bbbb"},
// 			File{Original: "aaaaa", Modified: "bbbbb"},
// 			File{Original: "aaaaaa", Modified: "bbbbbb"},
// 			File{Original: "aaaaaaa", Modified: "bbbbbbb"},
// 		},
// 	}
// 	want := " Path: ./some/dir\n" +
// 		" The following files will be renamed:\n" +
// 		"\n" +
// 		"    a       -> b\n" +
// 		"    aa      -> bb\n" +
// 		"    aaa     -> bbb\n" +
// 		"    aaaa    -> bbbb\n" +
// 		"    aaaaa   -> bbbbb\n" +
// 		"    aaaaaa  -> bbbbbb\n" +
// 		"    aaaaaaa -> bbbbbbb\n"

// 	var buff bytes.Buffer
// 	files.Display(&buff)

// 	got := buff.String()
// 	if want != got {
// 		t.Errorf("wanted:\n\n%s\n\ngot:\n\n%s", want, got)
// 	}
// }

// func TestTransform(t *testing.T) {
// 	testcases := []struct {
// 		name  string
// 		input string
// 		eplen int
// 		want  string
// 	}{
// 		{
// 			"tensorflow single digit",
// 			"1 - L1M1 Welcome To The Course Part 1 V1.mp4",
// 			2,
// 			"S01E01_L1M1_Welcome_To_The_Course_Part_1_V1.mp4",
// 		},
// 		{
// 			"tensorflow double digit",
// 			"10 - L2M1 Introduction V2.mp4",
// 			2,
// 			"S01E10_L2M1_Introduction_V2.mp4",
// 		},
// 		{
// 			"vuejs single digit",
// 			"#9 Creating and Using Components - VueJS For Everyone-2gpvyaaS1RI.mp4",
// 			2,
// 			"S01E09_Creating_and_Using_Components.mp4",
// 		},
// 		{
// 			"vuejs double digit",
// 			"#23 Lifecycle Methods - VueJS For Everyone-Ls8RGYKF68I.mp4",
// 			2,
// 			"S01E23_Lifecycle_Methods.mp4",
// 		},
// 		{
// 			"double digit out of 100",
// 			"99 - Espresso Outro.mp4",
// 			3,
// 			"S01E099_Espresso_Outro.mp4",
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			got := transform("01", tc.input, tc.eplen)
// 			if tc.want != got {
// 				t.Errorf("\ntransform(%q) returned:\n\t%q\nwanted:\n\t%q", tc.input, got, tc.want)
// 			}
// 		})
// 	}
// }
