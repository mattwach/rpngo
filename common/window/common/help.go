package common

import "mattwach/rpngo/common/rpn"

func RegisterConceptHelp(r *rpn.RPN, supportsLayout bool) {
	conceptHelp := map[string]string{
		"window.props": "Each window supports properties that changes how the window operates\n" +
			"- Print all properties and values for a window with w.listp\n" +
			"- Get a single property with w.getp\n" +
			"- Set a single property with w.setp\n" +
			"See Also: windows, window.layout, plotting",

		"windows": "The display can be customized with different windows\n" +
			"- Add a window with a w.new.<type> command. Example: w.new.stack\n" +
			"- Reset to a single window with w.reset.\n" +
			"See Also: window.layout, window.props",
	}

	if supportsLayout {
		conceptHelp["window.layout"] = "Windows are arranged with window groups.  There\n" +
			"is always a window group named 'root' which is the parent of all \n" +
			"windows and groups.\n" +
			"- Add a new window group to the root window with w.new.group.\n" +
			"- Move a window or group to a different window group with w.move.beg and w.move.end\n" +
			"- Change the weight of a window or group with w.weight (default weight is 100).\n" +
			"- Change the layout mode of a window group to columns with w.columns.\n" +
			"- Print info on all existing windows and groups with w.dump.\n" +
			"- You may also set .wtarget, .wend, and .wweight to direct how and\n" +
			"  where the next window/group will be create.  w.reset resets these\n" +
			"  to .wtarget=root, .wend=true, .wweight=100. Using illegal types or\n" +
			"  values for these variables will cause them to revert to the defaults\n" +
			"  as well.\n" +
			"See Also: windows, window.props"

	}
	r.RegisterConceptHelp(conceptHelp)
}

const WNewStackHelp = "Creates a new stack window with the given name and\n" +
	"adds it to the root window. Example: 's1' w.new.stack"

const WNewPlotHelp = "Creates a new plot window with the given name and\n" +
	"adds it to the root window. Example: 'p1' w.new.plot"
