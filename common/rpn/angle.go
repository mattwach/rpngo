package rpn

const setAngleHelp = "sets angle units to 'rad', 'deg', or 'grads'"

func setAngle(r *RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return ErrExpectedAString
	}
	switch f.UnsafeString() {
	case "rad":
		r.AngleUnit = POLAR_RAD_FRAME
	case "deg":
		r.AngleUnit = POLAR_DEG_FRAME
	case "grad":
		r.AngleUnit = POLAR_GRAD_FRAME
	default:
		return ErrChooseDegRadOGrad
	}
	return nil
}

const getAngleHelp = "returns currently-set angle units"

func getAngle(r *RPN) error {
	switch r.AngleUnit {
	case POLAR_RAD_FRAME:
		return r.PushFrame(StringFrame("rad", STRING_SINGLEQ_FRAME))
	case POLAR_DEG_FRAME:
		return r.PushFrame(StringFrame("deg", STRING_SINGLEQ_FRAME))
	case POLAR_GRAD_FRAME:
		return r.PushFrame(StringFrame("grad", STRING_SINGLEQ_FRAME))
	}
	return ErrIllegalValue
}

const radHelp = "sets trig / polar units to radians (calls 'rad' setangle)"

func rad(r *RPN) error {
	if err := r.PushFrame(StringFrame("rad", STRING_SINGLEQ_FRAME)); err != nil {
		return err
	}
	return setAngle(r)
}

const degHelp = "sets trig / polar units to degrees (calls 'deg' setangle)"

func deg(r *RPN) error {
	if err := r.PushFrame(StringFrame("deg", STRING_SINGLEQ_FRAME)); err != nil {
		return err
	}
	return setAngle(r)
}

const gradHelp = "sets trig / polar units to grads (calls 'grad' setangle)"

func grad(r *RPN) error {
	if err := r.PushFrame(StringFrame("grad", STRING_SINGLEQ_FRAME)); err != nil {
		return err
	}
	return setAngle(r)
}

func (r *RPN) FromRadians(rad complex128, t Frame) Frame {
	switch r.AngleUnit {
	case POLAR_DEG_FRAME:
		return Frame{ftype: t.ftype, cmplx: rad * 57.29577951308232, str: "`deg"}
	case POLAR_GRAD_FRAME:
		return Frame{ftype: t.ftype, cmplx: rad * 63.66197723675813, str: "`grad"}
	default:
		return Frame{ftype: t.ftype, cmplx: rad, str: "`rad"}
	}
}

func FromRadiansFloat(rad float64, t FrameType) float64 {
	switch t {
	case POLAR_DEG_FRAME:
		return rad * 57.29577951308232
	case POLAR_GRAD_FRAME:
		return rad * 63.66197723675813
	default:
		return rad
	}
}

func (r *RPN) ToRadians(angle complex128) complex128 {
	switch r.AngleUnit {
	case POLAR_DEG_FRAME:
		return angle * 0.0174532925199433
	case POLAR_GRAD_FRAME:
		return angle * 0.01570796326794897
	default:
		return angle
	}
}

func toRadiansFloat(angle float64, t FrameType) float64 {
	switch t {
	case POLAR_DEG_FRAME:
		return angle * 0.0174532925199433
	case POLAR_GRAD_FRAME:
		return angle * 0.01570796326794897
	default:
		return angle
	}
}
