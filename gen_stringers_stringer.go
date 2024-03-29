// Code generated by "stringer -type=GenAttr,Complete,Range,Event,NArgs,Flavor -linecomment -output gen_stringers_stringer.go"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AttrBang-0]
	_ = x[AttrBar-1]
	_ = x[AttrRegister-2]
	_ = x[AttrBuffer-3]
}

const _GenAttr_name = "-bang-bar-register-buffer"

var _GenAttr_index = [...]uint8{0, 5, 9, 18, 25}

func (i GenAttr) String() string {
	if i >= GenAttr(len(_GenAttr_index)-1) {
		return "GenAttr(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _GenAttr_name[_GenAttr_index[i]:_GenAttr_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CompleteArglist-0]
	_ = x[CompleteAugroup-1]
	_ = x[CompleteBuffer-2]
	_ = x[CompleteBehave-3]
	_ = x[CompleteColor-4]
	_ = x[CompleteCommand-5]
	_ = x[CompleteCompiler-6]
	_ = x[CompleteCscope-7]
	_ = x[CompleteDir-8]
	_ = x[CompleteEnvironment-9]
	_ = x[CompleteEvent-10]
	_ = x[CompleteExpression-11]
	_ = x[CompleteFile-12]
	_ = x[CompleteFileInPath-13]
	_ = x[CompleteFiletype-14]
	_ = x[CompleteFunction-15]
	_ = x[CompleteHelp-16]
	_ = x[CompleteHighlight-17]
	_ = x[CompleteHistory-18]
	_ = x[CompleteLocale-19]
	_ = x[CompleteMapclear-20]
	_ = x[CompleteMapping-21]
	_ = x[CompleteMenu-22]
	_ = x[CompleteMessages-23]
	_ = x[CompleteOption-24]
	_ = x[CompletePackadd-25]
	_ = x[CompleteShellCmd-26]
	_ = x[CompleteSign-27]
	_ = x[CompleteSyntax-28]
	_ = x[CompleteSyntime-29]
	_ = x[CompleteTag-30]
	_ = x[CompleteTagListFiles-31]
	_ = x[CompleteUser-32]
	_ = x[CompleteVar-33]
}

const _Complete_name = "-complete=arglist-complete=augroup-complete=buffer-complete=behave-complete=color-complete=command-complete=compiler-complete=cscope-complete=dir-complete=environment-complete=event-complete=expression-complete=file-complete=file_in_path-complete=filetype-complete=function-complete=help-complete=highlight-complete=history-complete=locale-complete=mapclear-complete=mapping-complete=menu-complete=messages-complete=option-complete=packadd-complete=shellcmd-complete=sign-complete=syntax-complete=syntime-complete=tag-complete=tag_listfiles-complete=user-complete=var"

var _Complete_index = [...]uint16{0, 17, 34, 50, 66, 81, 98, 116, 132, 145, 166, 181, 201, 215, 237, 255, 273, 287, 306, 323, 339, 357, 374, 388, 406, 422, 439, 457, 471, 487, 504, 517, 540, 554, 567}

func (i Complete) String() string {
	if i >= Complete(len(_Complete_index)-1) {
		return "Complete(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Complete_name[_Complete_index[i]:_Complete_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RangeLine-0]
	_ = x[RangeFile-1]
}

const _Range_name = "-range-range=%"

var _Range_index = [...]uint8{0, 6, 14}

func (i Range) String() string {
	if i >= Range(len(_Range_index)-1) {
		return "Range(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Range_name[_Range_index[i]:_Range_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NArgs0-0]
	_ = x[NArgs1-1]
	_ = x[NArgsZeroOrMore-2]
	_ = x[NArgsZeroOrOne-3]
	_ = x[NArgsOneOrMore-4]
}

const _NArgs_name = "-nargs=0-nargs=1-nargs=*-nargs=?-nargs=+"

var _NArgs_index = [...]uint8{0, 8, 16, 24, 32, 40}

func (i NArgs) String() string {
	if i >= NArgs(len(_NArgs_index)-1) {
		return "NArgs(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NArgs_name[_NArgs_index[i]:_NArgs_index[i+1]]
}

