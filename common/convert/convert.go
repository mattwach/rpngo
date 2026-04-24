package convert

import (
	"errors"
	"fmt"
	"math"
	"mattwach/rpngo/common/elog"
	"sort"
	"strings"
)

var (
	errUnknownConversionType = errors.New("unknown conversion type")
)

//
// Implementation:
//
// Step 1: interpret the type string via a series of aliases
// Step 2: flip term 1 if needed, or return an error if there is no way to
//   convert
// Step 3: convert numerators
// Step 4: If needed
//  a) invert
//  b) convert numerators
//  c) invert
//
// Done!

type unit struct {
	scale  float64
	offset float64
	names  []string
}

//
// Conversions to meters
//

var distantConvert = []unit{
	{1.0000, 0, []string{"meter", "m", "meters", "metre", "metres"}},
	{1609.3440, 0, []string{"mile", "miles", "mi"}},
	{0.0254, 0, []string{"inch", "inches", "in"}},
	{0.0100, 0, []string{"centimeter", "centimeters", "cm"}},
	{0.1000, 0, []string{"decimeter", "decimeters"}},
	{0.0010, 0, []string{"millimeter", "millimeters", "mm"}},
	{1000.0000, 0, []string{"kilometer", "kilometers", "km"}},
	{0.3048, 0, []string{"foot", "feet", "ft"}},
	{0.9144, 0, []string{"yard", "yards", "yd"}},
	{0.155849127913, 0, []string{"_gallonmeters"}},
	{0.0779245639837, 0, []string{"_pintmeter"}},
	{0.0170183442005, 0, []string{"_teaspoonmeters"}},
	{0.0245446996271, 0, []string{"_tablespoonmeters"}},
	{0.098178798467, 0, []string{"_quartmeter"}},
	{63.6149072152, 0, []string{"_acremeter"}},
}

//
// Conversion to seconds
//

var timeConvert = []unit{
	{1.0000, 0, []string{"second", "seconds", "sec", "s"}},
	{0.0010, 0, []string{"millisecond", "milliseconds", "ms"}},
	{1.0e-6, 0, []string{"microsecond", "microseconds", "us"}},
	{1.0e-9, 0, []string{"nanesecond", "nanoseconds", "ns"}},
	{60.0000, 0, []string{"minute", "minutes", "min"}},
	{3600.0000, 0, []string{"hour", "hours", "hr", "h"}},
	{86400.0000, 0, []string{"day", "days"}},
	{604800.0000, 0, []string{"week", "weeks"}},
}

//
// Weight/Mass Conversions (planet earth :) )
//

var massConvert = []unit{
	{1.00000, 0, []string{"gram", "grams", "g"}},
	{1000.00000, 0, []string{"kilogram", "kilograms", "kg"}},
	{28.3495231, 0, []string{"ounce", "ounces", "oz"}},
	{453.59237, 0, []string{"pound", "pounds", "lbs", "lb"}},
	{907184.74000, 0, []string{"ton", "tons"}},
	{101.9680, 0, []string{"newton", "newtons", "n"}},
	{101968.0000, 0, []string{"kilonewton", "kilonewtons"}},
	{10196800.0000, 0, []string{"_barnewton"}},
}

//
// Cycles
//

var cycleConvert = []unit{
	{1.0000, 0, []string{"cycle", "cycles"}},
	{1000.0000, 0, []string{"kilocycle", "kilocycles"}},
	{1000000.0000, 0, []string{"megacycle", "megacycles", "megahertz", "mhz"}},
	{1000000000.0000, 0, []string{"gigacycle", "gigacycles"}},
}

//
// Memory
//

var memoryConvert = []unit{
	{1.0000, 0, []string{"bit", "bits"}},
	{1024.0000, 0, []string{"kilobit", "kilobits"}},
	{8.0000, 0, []string{"byte", "bytes"}},
	{8192.0000, 0, []string{"kilobyte", "kilobytes", "kb"}},
	{8388608.0000, 0, []string{"megabyte", "megabytes", "mb"}},
	{8589934592.0000, 0, []string{"gigabyte", "gigabytes", "gb"}},
	{8796093022208.0000, 0, []string{"terabyte", "terabytes", "tb"}},
}

//
// Angles
//

var angleConvert = []unit{
	{1.0000, 0, []string{"radians", "rad"}},
	{math.Pi / 180.0, 0, []string{"degrees", "deg"}},
}

//
// Energy
//

var energyConvert = []unit{
	{1.0000, 0, []string{"joules", "j"}},
	{1000.0000, 0, []string{"kilojoules", "kj"}},
	{1000000.0000, 0, []string{"megajoules", "mj"}},
	{3600000.0000, 0, []string{"kwh"}},
	{745.699872, 0, []string{"hps"}},
	{1055.0600, 0, []string{"btu"}},        // EC standard
	{105506000.0000, 0, []string{"therm"}}, // EC standard
	{4.2000, 0, []string{"calorie", "cal", "calories"}},
	{4184.0000, 0, []string{"kilocalorie", "kilocalories"}},
}

//
// Temperature
//

var temperatureConvert = []unit{
	{1.0, 0, []string{"c", "celsius"}},
	{5.0 / 9.0, -32, []string{"f", "fahrenheit"}},
	{1.0, -273.15, []string{"k", "kelvin"}},
}

//
// Aliases to help things along
//

var aliases = map[string]string{
	"acre":        "_acremeter*_acremeter",
	"acres":       "acre",
	"bar":         "_barnewton/meter*meter",
	"cadence":     "cycles/minute",
	"gallon":      "_gallonmeters*_gallonmeters*_gallonmeters",
	"gallons":     "gallon",
	"ghz":         "gigacycles/second",
	"hp":          "hps/second",
	"hz":          "cycles/second",
	"khz":         "kilocycles/second",
	"kilowatt":    "kilojoules/second",
	"kilowatts":   "kilowatt",
	"kw":          "kilojoules/second",
	"liter":       "decimeter*decimeter*decimeter",
	"litre":       "liter",
	"liters":      "liter",
	"litres":      "liter",
	"milliliter":  "cm*cm*cm",
	"milliliters": "milliliter",
	"ml":          "milliliter",
	"megawatt":    "megajoules/second",
	"megawatts":   "megawatt",
	"mw":          "megajoules/second",
	"mhz":         "megacycles/second",
	"mph":         "miles/hour",
	"pascal":      "newton/meter*meter",
	"pascals":     "pascal",
	"pa":          "pascal",
	"kilopascal":  "kilonewtons/meter*meter",
	"kilopascals": "kilopascal",
	"kpa":         "kilopascal",
	"pint":        "_pintmeter*_pintmeter*_pintmeter",
	"pints":       "pint",
	"psi":         "pounds/inch*inch",
	"quart":       "_quartmeter*_quartmeter*_quartmeter",
	"quarts":      "quart",
	"rpm":         "cycles/minute",
	"tablespoon":  "_tablespoonmeters*_tablespoonmeters*_tablespoonmeters",
	"tablespoons": "tablespoon",
	"tsp":         "_teaspoonmeters*_teaspoonmeters*_teaspoonmeters",
	"teaspoon":    "tsp",
	"teaspoons":   "tsp",
	"watt":        "joules/second",
	"watts":       "watt",
}

type conversionType struct {
	className string
	scale     float64
	offset    float64
}

type Conversion struct {
	convertDict map[string]conversionType
}

type conversionData struct {
	numerator       []conversionType
	denominator     []conversionType
	inverted        bool
	numeratorName   string
	denominatorName string
}

func Init() *Conversion {
	elog.Heap("alloc: /convert/convert.go:214: c := &Conversion{convertDict: make(map[string]conversionType)}")
	c := &Conversion{convertDict: make(map[string]conversionType)} // object allocated on the heap: escapes at line 223
	c.insertKeys("Distance", distantConvert)
	c.insertKeys("Time", timeConvert)
	c.insertKeys("Force/Weight/Mass (Planet Earth)", massConvert)
	c.insertKeys("Cycles", cycleConvert)
	c.insertKeys("Memory", memoryConvert)
	c.insertKeys("Angles", angleConvert)
	c.insertKeys("Energy", energyConvert)
	c.insertKeys("Temperature", temperatureConvert)
	return c
}

func (c *Conversion) insertKeys(className string, data []unit) {
	for _, d := range data {
		for _, k := range d.names {
			if _, ok := c.convertDict[k]; ok {
				elog.Print("Error: duplicate conversion key: ", k)
			}
			c.convertDict[k] = conversionType{className, d.scale, d.offset}
		}
	}
}

func (c *Conversion) Convert(value float64, valueType string, targetType string) (float64, error) {
	// extract type and class information

	var source conversionData
	err := c.analyzeType(valueType, &source)
	if err != nil {
		return 0, err
	}

	var target conversionData
	err = c.analyzeType(targetType, &target)
	if err != nil {
		return 0, err
	}

	// check for ratio incompatibility

	if source.isRatio() != target.isRatio() {
		return 0, errors.New("can not convert between a ratio and scalar")
	}

	// check for inversion eligibility

	if target.isRatio() && (source.numeratorName == target.denominatorName) {
		target.invert()
	}

	// check for numerator compatibility

	if source.numeratorName != target.numeratorName {
		elog.Heap("alloc: /convert/convert.go:269: source.numeratorName,")
		return 0, fmt.Errorf(
			"incompatible numerator types: %s, %s",
			source.numeratorName, // object allocated on the heap: escapes at line 269
			target.numeratorName) // object allocated on the heap: escapes at line 270
	}

	// check for denominator compatibility, if needed

	if source.isRatio() && source.denominatorName != target.denominatorName {
		elog.Heap("alloc: /convert/convert.go:278: source.denominatorName,")
		return 0, fmt.Errorf(
			"incompatible denominator types: %s, %s",
			source.denominatorName, // object allocated on the heap: escapes at line 278
			target.denominatorName) // object allocated on the heap: escapes at line 279
	}

	// scale the value by each numerator

	for _, snum := range source.numerator {
		value = c.scaleUp(value, snum.scale, snum.offset)
	}

	for _, tnum := range target.numerator {
		value = c.scaleDown(value, tnum.scale, tnum.offset)
	}

	// If needed, scale by each denominator

	if source.isRatio() {
		for _, sden := range source.denominator {
			value = c.scaleDown(value, sden.scale, sden.offset)
		}
		for _, tden := range target.denominator {
			value = c.scaleUp(value, tden.scale, tden.offset)
		}
		if target.inverted {
			value = 1.0 / value
		}
	}

	return value, nil
}

func (c *Conversion) analyzeType(t string, data *conversionData) error {
	numeratorTypeList, denominatorTypeList := c.analyzeTypeStr(t)
	numeratorTypeList, denominatorTypeList = c.checkForAliases(numeratorTypeList, denominatorTypeList)
	denominatorTypeList, numeratorTypeList = c.checkForAliases(denominatorTypeList, numeratorTypeList)

	var numerator []conversionType
	var denominator []conversionType

	for _, numeratorType := range numeratorTypeList {
		n, ok := c.convertDict[numeratorType]
		if !ok {
			return fmt.Errorf("%v: %w", numeratorType, errUnknownConversionType)
		}
		numerator = append(numerator, n)
	}

	for _, denominatorType := range denominatorTypeList {
		d, ok := c.convertDict[denominatorType]
		if !ok {
			return fmt.Errorf("%v: %w", denominatorType, errUnknownConversionType)
		}
		denominator = append(denominator, d)
	}

	data.init(numerator, denominator)
	return nil
}

func (c *Conversion) analyzeTypeStr(t string) ([]string, []string) {
	var numeratorType []string
	var denominatorType []string
	parts := strings.SplitN(t, "/", 2)
	numeratorType = strings.Split(parts[0], "*")
	sort.Strings(numeratorType)
	if len(parts) > 1 {
		denominatorType = strings.Split(parts[1], "*")
		sort.Strings(denominatorType)
	}
	return numeratorType, denominatorType
}

func (c *Conversion) checkForAliases(numerator []string, denominator []string) ([]string, []string) {
	index := 0
	for index < len(numerator) {
		alias := aliases[numerator[index]]
		if alias != "" {
			numerator = append(numerator[:index], numerator[index+1:]...)
			parts := strings.SplitN(alias, "/", 2)
			numerator = append(numerator, strings.Split(parts[0], "*")...)
			if len(parts) > 1 {
				denominator = append(denominator, strings.Split(parts[1], "*")...)
			}
			index = 0
		} else {
			index++
		}
	}
	sort.Strings(numerator)
	sort.Strings(denominator)
	return numerator, denominator
}

func (c *conversionData) init(numerator, denominator []conversionType) {
	c.numerator = numerator
	c.denominator = denominator
	c.numeratorName = c.buildClassName(numerator, denominator)
	c.denominatorName = c.buildClassName(denominator, numerator)
}

func (cd *conversionData) buildClassName(numerator, denominator []conversionType) string {
	nameSet := make(map[string]bool)
	for _, n := range numerator {
		nameSet[n.className] = true
	}
	for _, d := range denominator {
		if nameSet[d.className] {
			nameSet[d.className] = false
		}
	}
	var nameList []string
	for _, n := range numerator {
		if nameSet[n.className] {
			nameList = append(nameList, n.className)
		}
	}
	return strings.Join(nameList, "*")
}

func (c *Conversion) scaleUp(value float64, scale float64, offset float64) float64 {
	return (value + offset) * scale
}

func (c *Conversion) scaleDown(value float64, scale float64, offset float64) float64 {
	return (value / scale) - offset
}

func (c *Conversion) Help() string {
	classes := make(map[string][]string)

	for name, conversion := range c.convertDict {
		classes[conversion.className] = append(classes[conversion.className], name)
	}

	var classNames []string
	for name := range classes {
		classNames = append(classNames, name)
	}
	sort.Strings(classNames)

	var lines []string
	for _, className := range classNames {
		lines = append(lines, "\n"+className+":")
		sort.Strings(classes[className])
		var subNames []string
		for _, name := range classes[className] {
			if (len(name) > 0) && (name[0] == '_') {
				continue
			}
			subNames = append(subNames, name)
			if len(subNames) == 4 {
				lines = dumpColumns(lines, subNames)
				subNames = subNames[:0]
			}
		}
		if len(subNames) > 0 {
			lines = dumpColumns(lines, subNames)
		}
	}

	lines = append(lines, "\nUseful Aliases:\n")
	var aliasNames []string
	for name := range aliases {
		aliasNames = append(aliasNames, name)
	}
	sort.Strings(aliasNames)

	for _, aliasName := range aliasNames {
		elog.Heap("alloc: convert/convert.go:451: fmt.Sprintf('  %-15s ->  %s\n', aliasName, aliases[aliasName]))")
		lines = append(
			lines,
			fmt.Sprintf("  %-15s ->  %s\n", aliasName, aliases[aliasName])) // object allocated on the heap: escapes at line 451
	}

	return strings.Join(lines, "")
}

func dumpColumns(lines []string, nameList []string) []string {
	lines = append(lines, "  ")
	for _, name := range nameList {
		lines = append(lines, fmt.Sprintf("%-15s ", name))
	}
	lines = append(lines, "\n")
	return lines
}

func (cd *conversionData) isRatio() bool {
	return len(cd.denominator) > 0
}

func (cd *conversionData) invert() {
	cd.numerator, cd.denominator = cd.denominator, cd.numerator
	cd.numeratorName, cd.denominatorName = cd.denominatorName, cd.numeratorName
	cd.inverted = !cd.inverted
}
