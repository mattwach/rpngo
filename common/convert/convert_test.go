package convert

import (
	"errors"
	"fmt"
	"testing"
)

func TestConvert(t *testing.T) {
	data := []struct{
		value float64
		valueType string
		targetType string
		// to work with float64 rounding, we multiply and convert to an int
		valMult float64
		wantVal int
		wantErr error
	}{
		// Distance
		{
			value: 5,
			valueType: "km",
			targetType: "mi",
			valMult: 1000,
			wantVal: 3106,
		},
		{
			value: 5,
			valueType: "km",
			wantErr: errUnknownConversionType,
		},
		{
			value: 5,
			targetType: "km",
			wantErr: errUnknownConversionType,
		},
		{
			value: 5,
			wantErr: errUnknownConversionType,
		},
		{
			value: 5,
			valueType: "inches",
			targetType: "cm",
			valMult: 10,
			wantVal: 127,
		},
		{
			value: 5,
			valueType: "inches",
			targetType: "decimeters",
			valMult: 100,
			wantVal: 127,
		},
		{
			value: 5,
			valueType: "ft",
			targetType: "mm",
			valMult: 1,
			wantVal: 1524,
		},
		{
			value: 5,
			valueType: "yard",
			targetType: "in",
			valMult: 1,
			wantVal: 180,
		},
		// Time
		{
			value: 5,
			valueType: "s",
			targetType: "ms",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "ms",
			targetType: "us",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "us",
			targetType: "ns",
			valMult: 1,
			wantVal: 4999, // rounding error
		},
		{
			value: 5,
			valueType: "min",
			targetType: "s",
			valMult: 1,
			wantVal: 300,
		},
		{
			value: 5,
			valueType: "hour",
			targetType: "min",
			valMult: 1,
			wantVal: 300,
		},
		{
			value: 5,
			valueType: "days",
			targetType: "hours",
			valMult: 1,
			wantVal: 120,
		},
		{
			value: 5,
			valueType: "weeks",
			targetType: "days",
			valMult: 1,
			wantVal: 35,
		},
		// Weight / mass
		{
			value: 5,
			valueType: "kg",
			targetType: "g",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "oz",
			targetType: "g",
			valMult: 10,
			wantVal: 1417,
		},
		{
			value: 5,
			valueType: "lb",
			targetType: "oz",
			valMult: 1,
			wantVal: 80,
		},
		{
			value: 5,
			valueType: "tons",
			targetType: "lb",
			valMult: 1,
			wantVal: 10000,
		},
		{
			value: 5,
			valueType: "lb",
			targetType: "newton",
			valMult: 100,
			wantVal: 2224,
		},
		{
			value: 5,
			valueType: "kilonewton",
			targetType: "newton",
			valMult: 1,
			wantVal: 5000,
		},
		// cycles
		{
			value: 5,
			valueType: "kilocycles",
			targetType: "cycles",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "megacycles",
			targetType: "kilocycles",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "gigacycles",
			targetType: "megacycles",
			valMult: 1,
			wantVal: 5000,
		},
		// memory
		{
			value: 5,
			valueType: "bytes",
			targetType: "bits",
			valMult: 1,
			wantVal: 40,
		},
		{
			value: 5,
			valueType: "kilobytes",
			targetType: "bytes",
			valMult: 1,
			wantVal: 1024 * 5,
		},
		{
			value: 5,
			valueType: "megabytes",
			targetType: "kilobytes",
			valMult: 1,
			wantVal: 1024 * 5,
		},
		{
			value: 5,
			valueType: "gigabytes",
			targetType: "megabytes",
			valMult: 1,
			wantVal: 1024 * 5,
		},
		{
			value: 5,
			valueType: "terabytes",
			targetType: "gigabytes",
			valMult: 1,
			wantVal: 1024 * 5,
		},
		// angles
		{
			value: 5,
			valueType: "radians",
			targetType: "degrees",
			valMult: 10,
			wantVal: 2864,
		},
		// energy
		{
			value: 5,
			valueType: "kj",
			targetType: "j",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "megajoules",
			targetType: "kilojoules",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "kwh",
			targetType: "kilojoules",
			valMult: 1,
			wantVal: 18000,
		},
		{
			value: 5,
			valueType: "hps",
			targetType: "joules",
			valMult: 1,
			wantVal: 3728,
		},
		{
			value: 5,
			valueType: "btu",
			targetType: "joules",
			valMult: 1,
			wantVal: 5275,
		},
		{
			value: 5,
			valueType: "therm",
			targetType: "btu",
			valMult: 1,
			wantVal: 500000,
		},
		{
			value: 5,
			valueType: "kilocalories",
			targetType: "joules",
			valMult: 1,
			wantVal: 20920,
		},
		// temperature
		{
			value: 5,
			valueType: "c",
			targetType: "f",
			valMult: 1,
			wantVal: 41,
		},
		{
			value: 5,
			valueType: "c",
			targetType: "k",
			valMult: 100,
			wantVal: 27814,
		},
		// aliases (and some expanded versions)
		{
			value: 5,
			valueType: "acre",
			targetType: "ft*ft",
			valMult: 1,
			wantVal: 217799, // rounding error
		},
		{
			value: 5,
			valueType: "bar",
			targetType: "psi",
			valMult: 100,
			wantVal: 7251,
		},
		{
			value: 5,
			valueType: "cadence",
			targetType: "cycles/s",
			valMult: 1000,
			wantVal: 83,
		},
		{
			value: 5,
			valueType: "gallon",
			targetType: "liter",
			valMult: 100,
			wantVal: 1892,
		},
		{
			value: 5,
			valueType: "ghz",
			targetType: "mhz",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "hp",
			targetType: "watts",
			valMult: 1,
			wantVal: 3728,
		},
		{
			value: 5,
			valueType: "khz",
			targetType: "hz",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "liter",
			targetType: "ml",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "mw",
			targetType: "kw",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "mi/h",
			targetType: "m/s",
			valMult: 1000,
			wantVal: 2235,
		},
		{
			value: 5,
			valueType: "mph",
			targetType: "m/s",
			valMult: 1000,
			wantVal: 2235,
		},
		{
			value: 5,
			valueType: "kpa",
			targetType: "pa",
			valMult: 1,
			wantVal: 5000,
		},
		{
			value: 5,
			valueType: "gallons",
			targetType: "pints",
			valMult: 1,
			wantVal: 39, // rounding error
		},
		{
			value: 5,
			valueType: "gallons",
			targetType: "quarts",
			valMult: 1,
			wantVal: 19, // rounding error
		},
		{
			value: 5,
			valueType: "rpm",
			targetType: "cadence",
			valMult: 1,
			wantVal: 5,
		},
		{
			value: 5,
			valueType: "pint",
			targetType: "tablespoon",
			valMult: 10,
			wantVal: 1599, // rounding error
		},
		{
			value: 5,
			valueType: "tablespoon",
			targetType: "teaspoon",
			valMult: 1,
			wantVal: 15, // rounding error
		},
		
		
	}

	for _, d := range data {
		name := fmt.Sprintf("%v %v>%v", d.value, d.valueType, d.targetType)
		t.Run(name, func(t *testing.T) {
			c := Init()
			val, err := c.Convert(d.value, d.valueType, d.targetType)
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err=%v, want %v", err, d.wantErr)
			}
			if d.wantErr != nil {
				return
			}
			if d.valMult == 0 {
				t.Errorf("d.valMult = 0 (test setup issue)")
			}
			gotVal := int(val * d.valMult)
			if d.wantVal != gotVal {
				t.Errorf("val=%v, want %v", gotVal, d.wantVal)
			}
		})
	}
}
