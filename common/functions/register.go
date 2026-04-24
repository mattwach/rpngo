package functions

import "mattwach/rpngo/common/rpn"

// RegisterAll resisters functions that are core to the RPN engine. There
// are other functions, such as window-specific functions, that will be
// added by their respective owners as the functions module should not
// know about these.
func RegisterAll(r *rpn.RPN) {
	r.Register("!=", notEqual, rpn.CatCompare, notEqualHelp)
	r.Register("<", lessThan, rpn.CatCompare, lessThanHelp)
	r.Register("<=", lessThanEqual, rpn.CatCompare, lessThanEqualHelp)
	r.Register("=", equal, rpn.CatCompare, equalHelp)
	r.Register(">", greaterThan, rpn.CatCompare, greaterThanHelp)
	r.Register(">=", greaterThanEqual, rpn.CatCompare, greaterThanEqualHelp)

	r.Register("&", and, rpn.CatBitwise, andHelp)
	r.Register("|", or, rpn.CatBitwise, orHelp)
	r.Register("^", xor, rpn.CatBitwise, xorHelp)
	r.Register("<<", shiftLeft, rpn.CatBitwise, shiftLeftHelp)
	r.Register(">>", shiftRight, rpn.CatBitwise, shiftRightHelp)

	r.Register("-", subtract, rpn.CatCore, subtractHelp)
	r.Register("*", multiply, rpn.CatCore, multiplyHelp)
	r.Register("/", divide, rpn.CatCore, divideHelp)
	r.Register("+", add, rpn.CatCore, addHelp)
	r.Register("%", mod, rpn.CatCore, modHelp)
	r.Register("false", falseFn, rpn.CatCore, falseHelp)
	r.Register("frac", frac, rpn.CatCore, fracHelp)
	r.Register("min", min, rpn.CatCore, minHelp)
	r.Register("max", max, rpn.CatCore, maxHelp)
	r.Register("neg", negate, rpn.CatCore, negateHelp)
	r.Register("round", round, rpn.CatCore, roundHelp)
	r.Register("true", trueFn, rpn.CatCore, trueHelp)
	r.Register("del", del, rpn.CatData, delHelp)
	r.Register("fields", fields, rpn.CatData, fieldsHelp)
	r.Register("filter", filter, rpn.CatData, filterHelp)
	r.Register("keep", keep, rpn.CatData, keepHelp)
	r.Register("reverse", reverse, rpn.CatData, reverseHelp)
	r.Register("sort", sortFn, rpn.CatData, sortHelp)

	r.Register("**", power, rpn.CatEng, powerHelp)
	r.Register("acos", acos, rpn.CatEng, acosHelp)
	r.Register("asin", asin, rpn.CatEng, asinHelp)
	r.Register("atan", atan, rpn.CatEng, atanHelp)
	r.Register("abs", abs, rpn.CatEng, absHelp)
	r.Register("cos", cos, rpn.CatEng, cosHelp)
	r.Register("log", log, rpn.CatEng, logHelp)
	r.Register("log10", log10, rpn.CatEng, log10Help)
	r.Register("rand", randFn, rpn.CatEng, randHelp)
	r.Register("sin", sin, rpn.CatEng, sinHelp)
	r.Register("sq", sq, rpn.CatEng, sqHelp)
	r.Register("sqrt", sqrt, rpn.CatEng, sqrtHelp)
	r.Register("tan", tan, rpn.CatEng, tanHelp)

	r.Register("hexdump", hexdump, rpn.CatIO, hexdumpHelp)
	r.Register("input", input, rpn.CatIO, inputHelp)
	r.Register("print", printFn, rpn.CatIO, printHelp)
	r.Register("printall", printall, rpn.CatIO, printallHelp)
	r.Register("println", printlnFn, rpn.CatIO, printlnHelp)
	r.Register("printlnx", printlnx, rpn.CatIO, printlnxHelp)
	r.Register("prints", prints, rpn.CatIO, printsHelp)
	r.Register("printsx", printsx, rpn.CatIO, printsxHelp)
	r.Register("printx", printx, rpn.CatIO, printxHelp)

	r.Register("@", exec, rpn.CatProg, execHelp)
	r.Register("delay", delay, rpn.CatProg, delayHelp)
	r.Register("error", errorFn, rpn.CatProg, errorHelp)
	r.Register("for", forFn, rpn.CatProg, forHelp)
	r.Register("if", If, rpn.CatProg, IfHelp)
	r.Register("ifelse", ifelse, rpn.CatProg, ifelseHelp)
	r.Register("noop", noop, rpn.CatProg, noopHelp)
	r.Register("time", timeFn, rpn.CatProg, timeHelp)
	r.Register("try", try, rpn.CatProg, tryHelp)

	r.Register("d", dropAll, rpn.CatStack, dropAllHelp)

	r.Register("heapstats", heapstats, rpn.CatStatus, heapstatsHelp)

	r.Register("bin", bin, rpn.CatType, binHelp)
	r.Register("float", floatFn, rpn.CatType, floatHelp)
	r.Register("hex", hex, rpn.CatType, hexHelp)
	r.Register("imag", imagFn, rpn.CatType, imagHelp)
	r.Register("int", intFn, rpn.CatType, intHelp)
	r.Register("oct", oct, rpn.CatType, octHelp)
	r.Register("phase", phase, rpn.CatType, phaseHelp)
	r.Register("polar", polar, rpn.CatType, polarHelp)
	r.Register("real", realFn, rpn.CatType, realHelp)
	r.Register("str", str, rpn.CatType, strHelp)
}
