package success

import (
	_ "embed"
)

var (
	//go:embed sample-1.pm
	sample1 string
	//go:embed result-1.txt
	result1 string
	//go:embed sample-2.pm
	sample2 string
	//go:embed result-2.txt
	result2 string
	//go:embed sample-3.pm
	sample3 string
	//go:embed result-3.txt
	result3 string
	//go:embed sample-4.pm
	sample4 string
	//go:embed result-4.txt
	result4 string
	//go:embed sample-5.pm
	sample5 string
	//go:embed result-5.txt
	result5 string
	//go:embed sample-6.pm
	sample6 string
	//go:embed result-6.txt
	result6 string
	//go:embed sample-7.pm
	sample7 string
	//go:embed result-7.txt
	result7 string
	//go:embed sample-8.pm
	sample8 string
	//go:embed result-8.txt
	result8 string
	//go:embed sample-9.pm
	sample9 string
	//go:embed result-9.txt
	result9 string
	//go:embed sample-10.pm
	sample10 string
	//go:embed result-10.txt
	result10 string
	//go:embed sample-11.pm
	sample11 string
	//go:embed result-11.txt
	result11 string
	//go:embed sample-12.pm
	sample12 string
	//go:embed result-12.txt
	result12 string
	//go:embed sample-13.pm
	sample13 string
	//go:embed result-13.txt
	result13 string
	//go:embed sample-14.pm
	sample14 string
	//go:embed result-14.txt
	result14 string
	//go:embed sample-15.pm
	sample15 string
	//go:embed result-15.txt
	result15 string
	//go:embed sample-16.pm
	sample16 string
	//go:embed result-16.txt
	result16 string
	//go:embed sample-17.pm
	sample17 string
	//go:embed result-17.txt
	result17 string
	//go:embed sample-18.pm
	sample18 string
	//go:embed result-18.txt
	result18 string
	//go:embed sample-19.pm
	sample19 string
	//go:embed result-19.txt
	result19 string
	//go:embed sample-20.pm
	sample20 string
	//go:embed result-20.txt
	result20 string
	//go:embed sample-21.pm
	sample21 string
	//go:embed result-21.txt
	result21 string
	//go:embed sample-22.pm
	sample22 string
	//go:embed result-22.txt
	result22 string
	//go:embed sample-23.pm
	sample23 string
	//go:embed result-23.txt
	result23 string
	//go:embed sample-24.pm
	sample24 string
	//go:embed result-24.txt
	result24 string
	//go:embed sample-25.pm
	sample25 string
	//go:embed result-25.txt
	result25 string
	//go:embed sample-26.pm
	sample26 string
	//go:embed result-26.txt
	result26 string
	//go:embed sample-27.pm
	sample27 string
	//go:embed result-27.txt
	result27 string
	//go:embed sample-28.pm
	sample28 string
	//go:embed result-28.txt
	result28 string
	//go:embed sample-29.pm
	sample29 string
	//go:embed result-29.txt
	result29 string
	//go:embed sample-30.pm
	sample30 string
	//go:embed result-30.txt
	result30 string
	//go:embed sample-31.pm
	sample31 string
	//go:embed result-31.txt
	result31 string
	//go:embed sample-32.pm
	sample32 string
	//go:embed result-32.txt
	result32 string
	//go:embed sample-33.pm
	sample33 string
	//go:embed result-33.txt
	result33 string
	//go:embed sample-34.pm
	sample34 string
	//go:embed result-34.txt
	result34 string
	//go:embed sample-35.pm
	sample35 string
	//go:embed result-35.txt
	result35 string
	//go:embed sample-36.pm
	sample36 string
	//go:embed result-36.txt
	result36 string
	//go:embed sample-37.pm
	sample37 string
	//go:embed result-37.txt
	result37 string
	//go:embed sample-38.pm
	sample38 string
	//go:embed result-38.txt
	result38 string
	//go:embed sample-39.pm
	sample39 string
	//go:embed result-39.txt
	result39 string
)

type Script struct {
	Code   string
	Result string
}

var Samples = map[string]Script{
	"sample-1.pm": {
		Code:   sample1,
		Result: result1,
	},
	"sample-2.pm": {
		Code:   sample2,
		Result: result2,
	},
	"sample-3.pm": {
		Code:   sample3,
		Result: result3,
	},
	"sample-4.pm": {
		Code:   sample4,
		Result: result4,
	},
	"sample-5.pm": {
		Code:   sample5,
		Result: result5,
	},
	"sample-6.pm": {
		Code:   sample6,
		Result: result6,
	},
	"sample-7.pm": {
		Code:   sample7,
		Result: result7,
	},
	"sample-8.pm": {
		Code:   sample8,
		Result: result8,
	},
	"sample-9.pm": {
		Code:   sample9,
		Result: result9,
	},
	"sample-10.pm": {
		Code:   sample10,
		Result: result10,
	},
	"sample-11.pm": {
		Code:   sample11,
		Result: result11,
	},
	"sample-12.pm": {
		Code:   sample12,
		Result: result12,
	},
	"sample-13.pm": {
		Code:   sample13,
		Result: result13,
	},
	"sample-14.pm": {
		Code:   sample14,
		Result: result14,
	},
	"sample-15.pm": {
		Code:   sample15,
		Result: result15,
	},
	"sample-16.pm": {
		Code:   sample16,
		Result: result16,
	},
	"sample-17.pm": {
		Code:   sample17,
		Result: result17,
	},
	"sample-18.pm": {
		Code:   sample18,
		Result: result18,
	},
	"sample-19.pm": {
		Code:   sample19,
		Result: result19,
	},
	"sample-20.pm": {
		Code:   sample20,
		Result: result20,
	},
	"sample-21.pm": {
		Code:   sample21,
		Result: result21,
	},
	"sample-22.pm": {
		Code:   sample22,
		Result: result22,
	},
	"sample-23.pm": {
		Code:   sample23,
		Result: result23,
	},
	"sample-24.pm": {
		Code:   sample24,
		Result: result24,
	},
	"sample-25.pm": {
		Code:   sample25,
		Result: result25,
	},
	"sample-26.pm": {
		Code:   sample26,
		Result: result26,
	},
	"sample-27.pm": {
		Code:   sample27,
		Result: result27,
	},
	"sample-28.pm": {
		Code:   sample28,
		Result: result28,
	},
	"sample-29.pm": {
		Code:   sample29,
		Result: result29,
	},
	"sample-30.pm": {
		Code:   sample30,
		Result: result30,
	},
	"sample-31.pm": {
		Code:   sample31,
		Result: result31,
	},
	"sample-32.pm": {
		Code:   sample32,
		Result: result32,
	},
	"sample-33.pm": {
		Code:   sample33,
		Result: result33,
	},
	"sample-34.pm": {
		Code:   sample34,
		Result: result34,
	},
	"sample-35.pm": {
		Code:   sample35,
		Result: result35,
	},
	"sample-36.pm": {
		Code:   sample36,
		Result: result36,
	},
	"sample-37.pm": {
		Code:   sample37,
		Result: result37,
	},
	"sample-38.pm": {
		Code:   sample38,
		Result: result38,
	},
	"sample-39.pm": {
		Code:   sample39,
		Result: result39,
	},
}
