package texx

import "github.com/qwerkon/interactive-yearly-planner/app/tex"

func EmphCell(text string) string {
	return tex.CellColor("black", tex.TextColor("white", text))
}
