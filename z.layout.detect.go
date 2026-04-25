package main

import (
	"net/http"

	. "klpm/lib/output"
)

func LayoutFromDeviceMode(mode string) string {
	if mode == deviceMobile { return layoutPhone }
	return layoutDesktop
}

func RequestLayout(r *http.Request) string {
	return LayoutFromDeviceMode(SessionDeviceMode(r))
}

func DeviceModeFromLayout(layout string) string {
	if layout == layoutPhone { return deviceMobile }
	return deviceDesktop
}

func DeviceConfirmHeadScript(mode string) string {
	x := deviceDesktop
	if mode0, ok := NormalizeDeviceMode(mode); ok { x = mode0 }
	return Str(
		`(function(){`,
		`var s="`, x, `";`,
		`var d=(window.innerWidth<768||window.matchMedia("(pointer:coarse)").matches)?"mobile":"desktop";`,
		`if(d===s){return;}`,
		`document.cookie="device="+d+"; path=/; max-age=31536000; samesite=lax";`,
		`if(document.cookie.indexOf("device="+d)>=0){location.replace(location.href);}`,
		`})();`,
	)
}
