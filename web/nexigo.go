package web

var NexigoControls = []*ControlHandler{
	NewControlHandler("Zoom, Absolute", []*Control{
		{Url: "/zoomin", Icon: "zoom_in", Multiplier: 1},
		{Url: "/zoomout", Icon: "zoom_out", Multiplier: -1},
	}),

	NewControlHandler("Pan, Absolute", []*Control{
		{Url: "/panleft", Icon: "arrow_back", Multiplier: -1},
		{Url: "/panright", Icon: "arrow_forward", Multiplier: 1},
	}),

	NewControlHandler("Tilt, Absolute", []*Control{
		{Url: "/tiltup", Icon: "arrow_upward", Multiplier: 1},
		{Url: "/tiltdown", Icon: "arrow_downward", Multiplier: -1},
	}),

	NewControlHandler("Brightness", []*Control{
		{Url: "/brightnessup", Icon: "brightness_high", Multiplier: 10},
		{Url: "/brightnessdown", Icon: "brightness_low", Multiplier: -10},
	}),

	NewControlHandler("Contrast", []*Control{
		{Url: "/contrastup", Icon: "contrast_square", Multiplier: 10},
		{Url: "/contrastdown", Icon: "exposure", Multiplier: -10},
	}),

	NewControlHandler("Saturation", []*Control{
		{Url: "/saturationup", Icon: "backlight_high", Multiplier: 10},
		{Url: "/saturationdown", Icon: "backlight_low", Multiplier: -10},
	}),

	// {Key: "Pan, Absolute", Url: "/panleft", Icon: "arrow_back", Multiplier: -1},
	// {Key: "Pan, Absolute", Url: "/panright", Icon: "arrow_forward", Multiplier: 1},
	// {Key: "Tilt, Absolute", Url: "/tiltup", Icon: "arrow_upward", Multiplier: -1},
	// {Key: "Tilt, Absolute", Url: "/tiltdown", Icon: "arrow_downward", Multiplier: 1},
	// {Key: "Brightness", Url: "/brightnessup", Icon: "brightness_high", Multiplier: 10},
	// {Key: "Brightness", Url: "/brightnessdown", Icon: "brightness_low", Multiplier: -10},
	// {Key: "Contrast", Url: "/contrastup", Icon: "contrast_square", Multiplier: 10},
	// {Key: "Contrast", Url: "/contrastdown", Icon: "exposure", Multiplier: -10},
	// {Key: "Saturation", Url: "/saturationup", Icon: "backlight_high", Multiplier: 10},
	// {Key: "Saturation", Url: "/saturationdown", Icon: "backlight_low", Multiplier: -10},
}
