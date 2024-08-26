package web

import "html/template"

func NexigoControlList(tmpl *template.Template) (handlers []*V4lHandler) {
	handlers = []*V4lHandler{
		NewV4lHandler("Zoom, Absolute", []*V4lControl{
			{url: "/zoomin", icon: "zoom_in", multiplier: 1},
			{url: "/zoomout", icon: "zoom_out", multiplier: -1},
		}, tmpl),

		NewV4lHandler("Pan, Absolute", []*V4lControl{
			{url: "/panleft", icon: "arrow_back", multiplier: -1},
			{url: "/panright", icon: "arrow_forward", multiplier: 1},
		}, tmpl),

		NewV4lHandler("Tilt, Absolute", []*V4lControl{
			{url: "/tiltup", icon: "arrow_upward", multiplier: 1},
			{url: "/tiltdown", icon: "arrow_downward", multiplier: -1},
		}, tmpl),

		NewV4lHandler("Brightness", []*V4lControl{
			{url: "/brightnessup", icon: "brightness_high", multiplier: 10},
			{url: "/brightnessdown", icon: "brightness_low", multiplier: -10},
		}, tmpl),

		NewV4lHandler("Contrast", []*V4lControl{
			{url: "/contrastup", icon: "contrast_square", multiplier: 10},
			{url: "/contrastdown", icon: "exposure", multiplier: -10},
		}, tmpl),

		NewV4lHandler("Saturation", []*V4lControl{
			{url: "/saturationup", icon: "backlight_high", multiplier: 10},
			{url: "/saturationdown", icon: "backlight_low", multiplier: -10},
		}, tmpl),

		// {Key: "Pan, Absolute", url: "/panleft", icon: "arrow_back", multiplier: -1},
		// {Key: "Pan, Absolute", url: "/panright", icon: "arrow_forward", multiplier: 1},
		// {Key: "Tilt, Absolute", url: "/tiltup", icon: "arrow_upward", multiplier: -1},
		// {Key: "Tilt, Absolute", url: "/tiltdown", icon: "arrow_downward", multiplier: 1},
		// {Key: "Brightness", url: "/brightnessup", icon: "brightness_high", multiplier: 10},
		// {Key: "Brightness", url: "/brightnessdown", icon: "brightness_low", multiplier: -10},
		// {Key: "Contrast", url: "/contrastup", icon: "contrast_square", multiplier: 10},
		// {Key: "Contrast", url: "/contrastdown", icon: "exposure", multiplier: -10},
		// {Key: "Saturation", url: "/saturationup", icon: "backlight_high", multiplier: 10},
		// {Key: "Saturation", url: "/saturationdown", icon: "backlight_low", multiplier: -10},
	}
	return
}
func NexigoMenuList(tmpl *template.Template) (handlers []*V4lHandler) {
	handlers = []*V4lHandler{
		NewV4lHandler("Zoom, Absolute", []*V4lControl{
			{url: "/zoommenu", icon: "search", multiplier: 1,
				controls: []*V4lControl{
					{url: "/zoomin", icon: "zoom_in", multiplier: 1},
					{url: "/zoomout", icon: "zoom_out", multiplier: -1},
				},
			}}, tmpl),
		NewV4lHandler("Pan, Absolute", []*V4lControl{
			{url: "/panmenu", icon: "arrows_outward", multiplier: 1,
				controls: []*V4lControl{
					{url: "/panleft", icon: "arrow_back", multiplier: -1},
					{url: "/panright", icon: "arrow_forward", multiplier: 1},
				},
			}}, tmpl),
		NewV4lHandler("Tilt, Absolute", []*V4lControl{
			{url: "/tiltmenu", icon: "height", multiplier: 1,
				controls: []*V4lControl{
					{url: "/tiltup", icon: "arrow_upward", multiplier: 1},
					{url: "/tiltdown", icon: "arrow_downward", multiplier: -1},
				},
			}}, tmpl),

		NewV4lHandler("Brightness", []*V4lControl{
			{url: "/brightnessup", icon: "brightness_high", multiplier: 10},
			{url: "/brightnessdown", icon: "brightness_low", multiplier: -10},
		}, tmpl),

		NewV4lHandler("Contrast", []*V4lControl{
			{url: "/contrastup", icon: "contrast_square", multiplier: 10},
			{url: "/contrastdown", icon: "exposure", multiplier: -10},
		}, tmpl),

		NewV4lHandler("Saturation", []*V4lControl{
			{url: "/saturationup", icon: "backlight_high", multiplier: 10},
			{url: "/saturationdown", icon: "backlight_low", multiplier: -10},
		}, tmpl),

		// {Key: "Pan, Absolute", url: "/panleft", icon: "arrow_back", multiplier: -1},
		// {Key: "Pan, Absolute", url: "/panright", icon: "arrow_forward", multiplier: 1},
		// {Key: "Tilt, Absolute", url: "/tiltup", icon: "arrow_upward", multiplier: -1},
		// {Key: "Tilt, Absolute", url: "/tiltdown", icon: "arrow_downward", multiplier: 1},
		// {Key: "Brightness", url: "/brightnessup", icon: "brightness_high", multiplier: 10},
		// {Key: "Brightness", url: "/brightnessdown", icon: "brightness_low", multiplier: -10},
		// {Key: "Contrast", url: "/contrastup", icon: "contrast_square", multiplier: 10},
		// {Key: "Contrast", url: "/contrastdown", icon: "exposure", multiplier: -10},
		// {Key: "Saturation", url: "/saturationup", icon: "backlight_high", multiplier: 10},
		// {Key: "Saturation", url: "/saturationdown", icon: "backlight_low", multiplier: -10},
	}
	return
}
