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
	"sample-8.pm": Script{
		Code:   sample8,
		Result: result8,
	},
	// "sample-9.pm":  Script{
	// 	Code:   sample9,
	// 	Result: result9,
	// },
	// "sample-10.pm": Script{
	// 	Code:   sample10,
	// 	Result: result10,
	// },
	// "sample-11.pm": Script{
	// 	Code:   sample11,
	// 	Result: result11,
	// },
	// "sample-12.pm": Script{
	// 	Code:   sample12,
	// 	Result: result12,
	// },
	// "sample-13.pm": Script{
	// 	Code:   sample13,
	// 	Result: result13,
	// },
	// "sample-14.pm": Script{
	// 	Code:   sample14,
	// 	Result: result14,
	// },
	// "sample-15.pm": Script{
	// 	Code:   sample15,
	// 	Result: result15,
	// },
	// "sample-16.pm": Script{
	// 	Code:   sample16,
	// 	Result: result16,
	// },
	// "sample-17.pm": Script{
	// 	Code:   sample17,
	// 	Result: result17,
	// },
	// "sample-18.pm": Script{
	// 	Code:   sample18,
	// 	Result: result18,
	// },
	// "sample-19.pm": Script{
	// 	Code:   sample19,
	// 	Result: result19,
	// },
	// "sample-20.pm": Script{
	// 	Code:   sample20,
	// 	Result: result20,
	// },
	// "sample-21.pm": Script{
	// 	Code:   sample21,
	// 	Result: result21,
	// },
	// "sample-22.pm": Script{
	// 	Code:   sample22,
	// 	Result: result22,
	// },
	// "sample-23.pm": Script{
	// 	Code:   sample23,
	// 	Result: result23,
	// },
	// "sample-24.pm": Script{
	// 	Code:   sample24,
	// 	Result: result24,
	// },
	// "sample-25.pm": Script{
	// 	Code:   sample25,
	// 	Result: result25,
	// },
	// "sample-26.pm": Script{
	// 	Code:   sample26,
	// 	Result: result26,
	// },
	// "sample-27.pm": Script{
	// 	Code:   sample27,
	// 	Result: result27,
	// },
	// "sample-28.pm": Script{
	// 	Code:   sample28,
	// 	Result: result28,
	// },
	// "sample-29.pm": Script{
	// 	Code:   sample29,
	// 	Result: result29,
	// },
	// "sample-30.pm": Script{
	// 	Code:   sample30,
	// 	Result: result30,
	// },
	// "sample-31.pm": Script{
	// 	Code:   sample31,
	// 	Result: result31,
	// },
	// "sample-32.pm": Script{
	// 	Code:   sample32,
	// 	Result: result32,
	// },
	// "sample-33.pm": Script{
	// 	Code:   sample33,
	// 	Result: result33,
	// },
	// "sample-34.pm": Script{
	// 	Code:   sample34,
	// 	Result: result34,
	// },
	// "sample-35.pm": Script{
	// 	Code:   sample35,
	// 	Result: result35,
	// },
	// "sample-36.pm": Script{
	// 	Code:   sample36,
	// 	Result: result36,
	// },
	// "sample-37.pm": Script{
	// 	Code:   sample37,
	// 	Result: result37,
	// },
	// "sample-38.pm": Script{
	// 	Code:   sample38,
	// 	Result: result38,
	// },
	// "sample-39.pm": Script{
	// 	Code:   sample39,
	// 	Result: result39,
	// },
}
